package jobtab

import (
	"encoding/json"
	"log"
	"net/rpc"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/config"
	"github.com/pja237/slurmcommander-dev/internal/defaults"
	"github.com/pja237/slurmcommander-dev/internal/scrpc"
)

var (
	cc config.ConfigContainer
)

func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

// Calls `squeue` to get job information for Jobs Tab
func GetSqueue(t time.Time) tea.Msg {

	var (
		sqJson SqueueJSON
		reply  scrpc.ReplyArgs
	)

	cmd := cc.Binpaths["squeue"]
	// if cc.SccCache is configured, call getSqueueFromCache()
	if cc.Sccache != "" {
		// connect and fetch using rpc
		client, err := rpc.DialHTTP("tcp", cc.Sccache)
		if err != nil {
			log.Fatal("dialing:", err)
		}

		// Synchronous call
		args := &scrpc.ReqArgs{
			Cid:  237,
			Cstr: defaults.AppName,
		}
		err = client.Call("SqueueCache.GetCachedSqueue", args, &reply)
		if err != nil {
			log.Fatal("RPC call ERROR:", err)
		}
		//fmt.Printf("GOT Reply no: %d len: %d\n", reply.CachedSqueue.Counter, len(reply.SqueueJSON.Jobs))
		sqJson = SqueueJSON(reply.SqueueJSON)

	} else {
		// Get the data by calling squeue cmd
		out, err := exec.Command(cmd, defaults.SqueueCmdSwitches...).CombinedOutput()
		if err != nil {
			return command.ErrorMsg{
				From:    "GetSqueue",
				ErrHelp: "Failed to run squeue command, check your scom.conf and set the correct paths there.",
				OrigErr: err,
			}
		}

		err = json.Unmarshal(out, &sqJson)
		if err != nil {
			return command.ErrorMsg{
				From:    "GetSqueue",
				ErrHelp: "squeue JSON failed to parse, note your slurm version and open an issue with us here: https://github.com/pja237/SlurmCommander-dev/issues/new/choose",
				OrigErr: err,
			}
		}
	}

	return sqJson
}

func TimedGetSqueue(l *log.Logger) tea.Cmd {
	l.Printf("TimedGetSqueue() start, tick: %d\n", cc.GetTick())
	return tea.Tick(cc.GetTick()*time.Second, GetSqueue)
}

func QuickGetSqueue(l *log.Logger) tea.Cmd {
	l.Printf("QuickGetSqueue() start\n")
	return tea.Tick(0*time.Second, GetSqueue)
}
