package clustertab

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
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
		log.Fatalf("Error exec sinfo: %s : %q\n", cmd, err)
	}

	err = json.Unmarshal(out, &siJson)
	if err != nil {
		log.Fatalf("Error unmarshall: %q\n", err)
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
