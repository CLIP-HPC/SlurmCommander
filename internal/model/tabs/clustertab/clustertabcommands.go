package clustertab

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/CLIP-HPC/SlurmCommander/internal/command"
	"github.com/CLIP-HPC/SlurmCommander/internal/config"
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
			ErrHelp: "sinfo JSON failed to parse, note your slurm version and open an issue with us here: https://github.com/CLIP-HPC/SlurmCommander/issues/new/choose",
			OrigErr: err,
		}
	}

	return siJson

}

func TimedGetSinfo(l *log.Logger) tea.Cmd {
	l.Printf("TimedGetSinfo() start, tick: %d\n", cc.GetTick())
	return tea.Tick(cc.GetTick()*time.Second, GetSinfo)
}

func QuickGetSinfo(l *log.Logger) tea.Cmd {
	l.Printf("QuickGetSinfo() start")
	return tea.Tick(0*time.Second, GetSinfo)
}
