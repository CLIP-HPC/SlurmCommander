package jobtab

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/styles"
)

type JobMenuOptions map[string]MenuOptions

type MenuOptions []list.Item
type MenuItem struct {
	action      string
	description string
}

var MenuList = JobMenuOptions{
	"PENDING": MenuOptions{
		MenuItem{
			action:      "INFO",
			description: "Show job information",
		},
		MenuItem{
			action:      "CANCEL",
			description: "Cancel the selected job",
		},
		MenuItem{
			action:      "HOLD",
			description: "Prevent a job from starting",
		},
	},
	"RUNNING": MenuOptions{
		MenuItem{
			action:      "INFO",
			description: "Show job information",
		},
		MenuItem{
			action:      "CANCEL",
			description: "Cancel the selected job",
		},
		MenuItem{
			action:      "SSH",
			description: "Ssh connect to job batch node",
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

func (m *MenuItem) ExecMenuItem(jobID string, node string, l *log.Logger) tea.Cmd {

	l.Printf("ExecMenuItem() jobID=%s node=%s m.action=%s\n", jobID, node, m.action)

	// TODO: move related commands to jobtab package
	switch m.action {
	case "CANCEL":
		return command.CallScancel(jobID, l)
	case "HOLD":
		return command.CallScontrolHold(jobID, l)
	case "REQUEUE":
		return command.CallScontrolRequeue(jobID, l)
	case "SSH":
		return command.CallSsh(node, l)
	}

	return nil
}

func NewMenu(selJobState string, l *log.Logger) list.Model {
	var lm list.Model = list.Model{}

	menuOpts := MenuList[selJobState]
	l.Printf("MENU Options%#v\n", MenuList[selJobState])

	defDel := list.NewDefaultDelegate()
	defStyles := list.NewDefaultItemStyles()
	defStyles.NormalTitle = styles.MenuNormalTitle
	defStyles.SelectedTitle = styles.MenuSelectedTitle
	//defStyles.NormalDesc = styles.MenuNormalDesc
	defStyles.SelectedDesc = styles.MenuSelectedDesc
	defDel.Styles = defStyles
	lm = list.New(menuOpts, defDel, 10, 10)
	lm.Title = "Job actions"
	lm.SetShowStatusBar(true)
	lm.SetFilteringEnabled(false)
	lm.SetShowHelp(false)
	lm.SetShowPagination(false)
	lm.SetSize(30, 25)
	lm.Styles.Title = styles.MenuTitleStyle

	return lm
}
