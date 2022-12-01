package main

import (
	"encoding/json"
	"flag"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os/exec"
	"strconv"
	"time"

	"github.com/pja237/slurmcommander-dev/internal/defaults"
	"github.com/pja237/slurmcommander-dev/internal/scrpc"
)

type cmdflags struct {
	port    *uint
	tick    *uint
	prefix  *string
	version *bool
}

func (c *cmdflags) GetCmdFlags() {
	c.port = flag.Uint("p", defaults.SccPort, "Port where to listen for SC requests")
	c.tick = flag.Uint("t", defaults.SccRefreshT, "Time period, how often to fetch data from slurm")
	c.prefix = flag.String("f", defaults.SccPrefix, "preFix where slurm commands are found")
	c.version = flag.Bool("v", false, "show version")
	flag.Parse()

}

func (c *cmdflags) DumpFlags() {
	log.Println("--------------------------------------------------------------------------------")
	log.Printf("Port: %d\n", *c.port)
	log.Printf("Version: %t\n", *c.version)
	log.Printf("Tick: %d\n", *c.tick)
	log.Println("--------------------------------------------------------------------------------")
}

// listener->bank request data
type lbQuery struct {
	req     scrpc.ReqArgs
	replyCh chan scrpc.CachedSqueue
}

// Data type (receiver) whose methods are published in RPC server
type SqueueCache struct {
	lbCh chan lbQuery
}

func (sqc *SqueueCache) GetCachedSqueue(c scrpc.ReqArgs, r *scrpc.ReplyArgs) error {
	var lbq lbQuery

	t := time.Now()
	log.Printf(">> GetCachedSqueue() invoked")
	log.Printf(">> GetCachedSqueue() client id=%d str=%s", c.Cid, c.Cstr)

	// Send a query to the bank
	log.Printf(">> GetCachedSqueue() Prep req for bank, t=%f\n", time.Since(t).Seconds())
	lbq.replyCh = make(chan scrpc.CachedSqueue)
	lbq.req.Cid = c.Cid
	lbq.req.Cstr = c.Cstr
	sqc.lbCh <- lbq

	// get response
	log.Printf(">> GetCachedSqueue() Send req to bank, t=%f\n", time.Since(t).Seconds())
	csq := <-lbq.replyCh
	log.Printf(">> GetCachedSqueue() Got response from bank, len(resp)=%d\n", len(csq.SqueueJSON.Jobs))

	// fill out ReplyArgs
	r.N = len(csq.SqueueJSON.Jobs)
	r.CachedSqueue = csq

	log.Printf(">> GetCachedSqueue() Done: t=%fs", time.Since(t).Seconds())
	return nil
}

type routineDone struct {
	err error
}

func fetcherSacct(prefix string, fch chan<- routineDone, fbCh chan<- scrpc.CachedSqueue) {
	var (
		cSqueue scrpc.CachedSqueue = scrpc.CachedSqueue{}
	)

	defer func() {
		fch <- routineDone{
			err: nil,
		}
	}()

	// TODO: either modify and reuse GetSqueue from jobtabcommands.go or redo it here

	// ticker loop {
	for {

		// exec command and fetch data
		//out, err := exec.Command(prefix+"/squeue", defaults.SqueueCmdSwitches...).CombinedOutput()
		out, err := exec.Command(prefix, defaults.SqueueCmdSwitches...).CombinedOutput()
		if err != nil {
			// TODO: signalize error... to someone... something
			log.Printf("> fetcher: ERROR exec(): %s\n", err)
			break
		}

		err = json.Unmarshal(out, &cSqueue.SqueueJSON)
		if err != nil {
			// TODO: signalize error... to someone... something
			log.Printf("> fetcher: ERROR unmarshall(): %s\n", err)
			break
		}
		// all went well, increment the counter
		cSqueue.Counter++
		log.Printf("> fetcher: counter=%d\n", cSqueue.Counter)
		log.Printf("> fetcher: len([]jobs)=%d\n", len(cSqueue.SqueueJSON.Jobs))
		fbCh <- cSqueue

		// send via channel

		time.Sleep(5 * time.Second)
	}

	// } eoticker loop

}

func bank(bCh chan<- routineDone, fbCh <-chan scrpc.CachedSqueue, lbCh <-chan lbQuery) {

	var (
		csq scrpc.CachedSqueue
	)

	defer func() {
		bCh <- routineDone{
			err: nil,
		}
	}()

	ticker := time.Tick(1 * time.Second)

	for {
		select {
		case csq = <-fbCh:
			log.Printf("> bank: got msg no. %d from fetcher, len([]jobs)=%d\n", csq.Counter, len(csq.SqueueJSON.Jobs))
		case lbReq := <-lbCh:
			log.Printf("> bank: got req from listener id=%d str=%s\n", lbReq.req.Cid, lbReq.req.Cstr)
			log.Printf("> bank: sending reply to listener")
			// TODO: what if bank has nothing yet?
			lbReq.replyCh <- csq
		case <-ticker:
			log.Printf("> bank: tick\n")

		}

		//time.Sleep(1 * time.Second)
	}

}

func listenRPC(p uint, lbCh chan lbQuery) {

	sqc := new(SqueueCache)
	sqc.lbCh = lbCh

	rpc.Register(sqc)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":"+strconv.FormatUint(uint64(p), 10))
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
	log.Printf("> listenRPC setup done")
}

func main() {

	var (
		cmd  cmdflags
		fCh  chan routineDone        // Fetcher done channel
		bCh  chan routineDone        // Bank done channel
		fbCh chan scrpc.CachedSqueue // Fetcher-Bank channel
		lbCh chan lbQuery            // ListenRPC-Bank channel
	)

	fCh = make(chan routineDone)
	bCh = make(chan routineDone)
	fbCh = make(chan scrpc.CachedSqueue)
	lbCh = make(chan lbQuery, 10)

	log.Printf("- SCCache Start")

	// Parse CMDline parameters
	cmd.GetCmdFlags()
	cmd.DumpFlags()

	// Read config
	//  - skip, flags only

	// Setup logger
	//	- skip

	// Spin up scraper goroutine
	go fetcherSacct(*cmd.prefix, fCh, fbCh)

	// Spin up bank
	go bank(bCh, fbCh, lbCh)

	// Spin up server listeners
	go listenRPC(*cmd.port, lbCh)

	log.Printf("Spun up goroutines, waiting...\n")
	// TODO: select on goroutine exit channels, if one exits, log, teardown everything and exit
	select {
	case err := <-fCh:
		if err.err != nil {
			log.Printf("Got ERROR from fetcher: %s\n", err.err)
		} else {
			log.Printf("Fetcher finished: OK")
		}
	case err := <-bCh:
		if err.err != nil {
			log.Printf("Got ERROR from bank: %s\n", err.err)
		} else {
			log.Printf("Bank finished: OK")
		}
	}

	log.Printf("- SCCache End\n")
}
