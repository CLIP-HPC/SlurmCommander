package jobhisttab

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander/internal/keybindings"
	"github.com/pja237/slurmcommander/internal/slurm"
)

type JobHistTab struct {
	SacctTable        table.Model
	SacctHist         slurm.SacctJobHist
	SacctHistFiltered slurm.SacctJobHist
	Filter            textinput.Model
	Stats
}

type Stats struct {
	StateCnt map[string]uint
	AvgWait  time.Duration
	SDWait   int
}

func (t *JobHistTab) GetStatsFiltered(l *log.Logger) {
	t.Stats.StateCnt = map[string]uint{}
	t.AvgWait = 0

	l.Printf("GetStatsFiltered start\n")
	for _, v := range t.SacctHistFiltered.Jobs {
		t.Stats.StateCnt[*v.State.Current]++
		t.AvgWait += time.Unix(int64(*v.Time.Start), 0).Sub(time.Unix(int64(*v.Time.Submission), 0))
	}
	l.Printf("GetStatsFiltered totalwait: %d\n", t.AvgWait)
	l.Printf("GetStatsFiltered totalwait: %s\n", t.AvgWait.String())
	t.AvgWait = t.AvgWait / time.Duration((len(t.SacctHistFiltered.Jobs)))
	l.Printf("GetStatsFiltered avgwait: %d\n", t.AvgWait)
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
	&keybindings.DefaultKeyMap.Enter:    true,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
}
