package jobtab

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/keybindings"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
	"github.com/pja237/slurmcommander-dev/internal/stats"
	"github.com/pja237/slurmcommander-dev/internal/styles"
)

type JobTab struct {
	InfoOn           bool
	StatsOn          bool
	SqueueTable      table.Model
	Squeue           slurm.SqueueJSON
	SqueueFiltered   slurm.SqueueJSON
	Filter           textinput.Model
	SelectedJob      string
	SelectedJobState string
	MenuOn           bool
	MenuChoice       MenuItem
	Menu             list.Model
	Stats
}

type Stats struct {
	// TODO: also perhaps: count by user? account?
	StateCnt map[string]uint
	AvgWait  time.Duration
	MinWait  time.Duration
	MaxWait  time.Duration
	MedWait  time.Duration
	AvgRun   time.Duration
	MinRun   time.Duration
	MaxRun   time.Duration
	MedRun   time.Duration
}

type JobMenuOptions map[string]MenuOptions

type MenuOptions []list.Item
type MenuItem struct {
	action      string
	description string
}

func (t *JobTab) GetStatsFiltered(l *log.Logger) {
	t.Stats.StateCnt = map[string]uint{}
	tmp := []time.Duration{}
	tmpRun := []time.Duration{}
	t.AvgWait = 0
	t.MedWait = 0

	l.Printf("jobtab GetStatsFiltered start\n")
	for _, v := range t.SqueueFiltered.Jobs {
		t.Stats.StateCnt[*v.JobState]++
		switch *v.JobState {
		case "PENDING":
			tmp = append(tmp, time.Since(time.Unix(int64(*v.SubmitTime), 0)))
		case "RUNNING":
			tmpRun = append(tmp, time.Since(time.Unix(int64(*v.StartTime), 0)))
		}
	}
	l.Printf("jobtab GetStatsFiltered totalwait: %d\n", t.AvgWait)
	l.Printf("jobtab GetStatsFiltered totalwait: %s\n", t.AvgWait.String())
	t.MedWait, t.MinWait, t.MaxWait = stats.Median(tmp)
	t.MedRun, t.MinRun, t.MaxRun = stats.Median(tmpRun)
	t.AvgWait = stats.Avg(tmp)
	t.AvgRun = stats.Avg(tmpRun)
	l.Printf("jobtab GetStatsFiltered avgwait: %d\n", t.AvgWait)
	l.Printf("jobtab GetStatsFiltered medwait: %d\n", t.MedWait)
	l.Printf("jobtab GetStatsFiltered end\n")
	l.Printf("jobtab GetStatsFiltered end\n")
}

// TODO: we don't need to return messages, we're called from update, just error and let update continue...
func (m *MenuItem) ExecMenuItem(jobID string, l *log.Logger) tea.Cmd {
	//var msg tea.Msg

	l.Printf("ExecMenuItem() jobID=%s m.action=%s\n", jobID, m.action)

	switch m.action {
	case "CANCEL":
		// TODO: here call cancel method, return jobCancelMsg
		return command.CallScancel(jobID, l)
	case "HOLD":
		return command.CallScontrolHold(jobID, l)
	case "REQUEUE":
		return command.CallScontrolRequeue(jobID, l)
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
	&keybindings.DefaultKeyMap.Stats:    true,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
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
	lm.SetHeight(30)
	lm.SetWidth(30)
	lm.SetSize(30, 30)
	lm.Styles.Title = styles.MenuTitleStyle

	return lm
}
