package jobhisttab

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander-dev/internal/keybindings"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
	"github.com/pja237/slurmcommander-dev/internal/stats"
)

type JobHistTab struct {
	SacctTable        table.Model
	SacctHist         slurm.SacctJobHist
	SacctHistFiltered slurm.SacctJobHist
	Filter            textinput.Model
	Stats
}

type Stats struct {
	StatsOn  bool
	StateCnt map[string]uint
	AvgWait  time.Duration
	MinWait  time.Duration
	MaxWait  time.Duration
	MedWait  time.Duration
	SDWait   int
}

func (t *JobHistTab) GetStatsFiltered(l *log.Logger) {
	t.Stats.StateCnt = map[string]uint{}
	tmp := []time.Duration{}
	t.AvgWait = 0
	t.MedWait = 0

	l.Printf("GetStatsFiltered start on %d rows\n", len(t.SacctHistFiltered.Jobs))

	for _, v := range t.SacctHistFiltered.Jobs {
		t.Stats.StateCnt[*v.State.Current]++
		tmp = append(tmp, time.Unix(int64(*v.Time.Start), 0).Sub(time.Unix(int64(*v.Time.Submission), 0)))
	}

	l.Printf("GetStatsFiltered totalwait: %d\n", t.AvgWait)
	l.Printf("GetStatsFiltered totalwait: %s\n", t.AvgWait.String())
	t.MedWait, t.MinWait, t.MaxWait = stats.Median(tmp)
	t.AvgWait = stats.Avg(tmp)
	l.Printf("GetStatsFiltered avgwait: %d\n", t.AvgWait)
	l.Printf("GetStatsFiltered medwait: %d\n", t.MedWait)
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
	&keybindings.DefaultKeyMap.Stats:    true,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
}
