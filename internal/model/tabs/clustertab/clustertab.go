package clustertab

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander-dev/internal/keybindings"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
)

type JobClusterTab struct {
	SinfoTable    table.Model
	CpuBar        progress.Model
	MemBar        progress.Model
	Sinfo         slurm.SinfoJSON
	SinfoFiltered slurm.SinfoJSON
	Filter        textinput.Model
	Stats
}

type Stats struct {
	// TODO: also perhaps: count by user? account?
	StateCnt map[string]uint
}

func (t *JobClusterTab) GetStatsFiltered(l *log.Logger) {
	var key string

	t.Stats.StateCnt = map[string]uint{}

	l.Printf("GetStatsFiltered JobClusterTab start\n")
	for _, v := range t.SinfoFiltered.Nodes {
		if len(*v.StateFlags) != 0 {
			key = *v.State + "+" + strings.Join(*v.StateFlags, "+")
		} else {
			key = *v.State
		}
		//t.Stats.StateCnt[*v.JobState]++
		t.Stats.StateCnt[key]++
	}
	l.Printf("GetStatsFiltered end\n")
}

type Keys map[*key.Binding]bool

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.Up:       true,
	&keybindings.DefaultKeyMap.Down:     true,
	&keybindings.DefaultKeyMap.PageUp:   true,
	&keybindings.DefaultKeyMap.PageDown: true,
	&keybindings.DefaultKeyMap.Slash:    true,
	&keybindings.DefaultKeyMap.Info:     false,
	&keybindings.DefaultKeyMap.Enter:    false,
	&keybindings.DefaultKeyMap.Stats:    true,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
}
