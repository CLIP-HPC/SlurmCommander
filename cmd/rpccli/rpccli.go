package main

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/pja237/slurmcommander-dev/internal/slurm"
)

type CachedSqueue struct {
	Counter    uint
	SqueueJSON slurm.SqueueJSON
}

type ReqArgs struct {
	Cid  uint   // client id
	Cstr string // client string
}

type ReplyArgs struct {
	N int
	CachedSqueue
}

func main() {
	var reply ReplyArgs

	client, err := rpc.DialHTTP("tcp", "localhost:1234")
	if err != nil {
		log.Fatal("dialing:", err)
	}

	// Synchronous call
	args := &ReqArgs{237, "pja client"}
	err = client.Call("SqueueCache.GetCachedSqueue", args, &reply)
	if err != nil {
		log.Fatal("RPC call ERROR:", err)
	}
	fmt.Printf("GOT Reply no: %d len: %d\n", reply.CachedSqueue.Counter, len(reply.SqueueJSON.Jobs))
}
