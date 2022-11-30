package jobtab

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/config"
	"github.com/pja237/slurmcommander-dev/internal/defaults"
)

var (
	cc config.ConfigContainer
)

func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

// Calls `squeue` to get job information for Jobs Tab
func GetSqueue(t time.Time) tea.Msg {

	var sqJson SqueueJSON

	cmd := cc.Binpaths["squeue"]
	// if cc.SccCache is configured, call getSqueueFromCache()
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
