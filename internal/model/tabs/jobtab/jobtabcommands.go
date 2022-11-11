package jobtab

import (
	"encoding/json"
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/config"
)

var (
	cc                config.ConfigContainer
	SqueueCmdSwitches = []string{"-a", "--json"}
)

func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

// Calls `squeue` to get job information for Jobs Tab
func GetSqueue(t time.Time) tea.Msg {

	var sqJson SqueueJSON

	cmd := cc.Binpaths["squeue"]
	out, err := exec.Command(cmd, SqueueCmdSwitches...).CombinedOutput()
	if err != nil {
		log.Fatalf("Error exec squeue: %s : %q\n", cmd, err)
	}

	err = json.Unmarshal(out, &sqJson)
	if err != nil {
		log.Fatalf("Error unmarshall: %q\n", err)
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
