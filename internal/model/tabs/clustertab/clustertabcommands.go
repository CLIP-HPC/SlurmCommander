package clustertab

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/config"
)

var (
	cc               config.ConfigContainer
	SinfoCmdSwitches = []string{"-a", "--json"}
)

func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

// Calls `sinfo` to get node information for Cluster Tab
func GetSinfo(t time.Time) tea.Msg {
	var siJson SinfoJSON

	cmd := cc.Binpaths["sinfo"]
	out, err := exec.Command(cmd, SinfoCmdSwitches...).CombinedOutput()
	if err != nil {
		return command.ErrorMsg{
			From:    "GetSinfo",
			ErrHelp: "Failed to run sinfo command, check your scom.conf and set the correct paths there.",
			OrigErr: err,
		}
	}

	err = json.Unmarshal(out, &siJson)
	if err != nil {
		return command.ErrorMsg{
			From:    "GetSinfo",
			ErrHelp: "sinfo JSON failed to parse, note your slurm version and open an issue with us here: https://github.com/pja237/SlurmCommander-dev/issues/new/choose",
			OrigErr: err,
		}
	}

	return siJson

}

func TimedGetSinfo() tea.Cmd {
	// TODO: make timers configurable
	return tea.Tick(cc.GetTick()*time.Second, GetSinfo)
}

func QuickGetSinfo(l *log.Logger) tea.Cmd {
	l.Printf("QuickGetSinfo() start")
	return tea.Tick(0*time.Second, GetSinfo)
}
