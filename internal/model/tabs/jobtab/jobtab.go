package jobtab

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander/internal/keybindings"
	"github.com/pja237/slurmcommander/internal/slurm"
)

type JobTab struct {
	InfoOn           bool
	SqueueTable      table.Model
	Squeue           slurm.SqueueJSON
	SqueueFiltered   slurm.SqueueJSON
	Filter           textinput.Model
	SelectedJob      string
	SelectedJobState string
	MenuOn           bool
	MenuChoice       MenuItem
	Menu             list.Model
}

type JobMenuOptions map[string]MenuOptions

type MenuOptions []list.Item
type MenuItem struct {
	action      string
	description string
}

type jobCancelMsg struct{}
type jobHoldMsg struct{}

// TODO: we don't need to return messages, we're called from update, just error and let update continue...
func (m *MenuItem) ExecMenuItem(jobID string, l *log.Logger) tea.Msg {
	//var msg tea.Msg

	l.Printf("ExecMenuItem() jobID=%s m.action=%s\n", jobID, m.action)

	switch m.action {
	case "CANCEL":
		// TODO: here call cancel method, return jobCancelMsg
	case "HOLD":
	case "REQUEUE":
	}

	return nil
}

var MenuList = JobMenuOptions{
	"PENDING": MenuOptions{
		MenuItem{
			action:      "INFO",
			description: "Show job information",
		},
		MenuItem{
			action:      "CANCEL",
			description: "Cancels the selected job",
		},
		MenuItem{
			action:      "HOLD",
			description: "Lower job priority to 0.",
		},
	},
	"RUNNING": MenuOptions{
		MenuItem{
			action:      "INFO",
			description: "Show job information",
		},
		MenuItem{
			action:      "CANCEL",
			description: "Cancels the selected job",
		},
		MenuItem{
			action:      "REQUEUE",
			description: "Stop the job and send it back to Queue",
		},
	},
}

func (i MenuItem) FilterValue() string {
	return ""
}
func (i MenuItem) Title() string {
	return i.action
}
func (i MenuItem) Description() string {
	return i.description
}

type Keys map[*key.Binding]bool

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.Up:       true,
	&keybindings.DefaultKeyMap.Down:     true,
	&keybindings.DefaultKeyMap.PageUp:   true,
	&keybindings.DefaultKeyMap.PageDown: true,
	&keybindings.DefaultKeyMap.Slash:    true,
	&keybindings.DefaultKeyMap.Info:     true,
	&keybindings.DefaultKeyMap.Enter:    true,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
}
