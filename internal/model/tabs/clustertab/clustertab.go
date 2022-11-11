package clustertab

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

type JobClusterTab struct {
	StatsOn       bool
	SinfoTable    table.Model
	CpuBar        progress.Model
	MemBar        progress.Model
	Sinfo         SinfoJSON
	SinfoFiltered SinfoJSON
	Filter        textinput.Model
	Stats
}

type Stats struct {
	// TODO: also perhaps: count by user? account?
	StateCnt       map[string]uint
	StateSimpleCnt map[string]uint
}

func (t *JobClusterTab) GetStatsFiltered(l *log.Logger) {
	var key string

	t.Stats.StateCnt = map[string]uint{}
	t.Stats.StateSimpleCnt = map[string]uint{}

	l.Printf("GetStatsFiltered JobClusterTab start\n")
	for _, v := range t.SinfoFiltered.Nodes {
		if len(*v.StateFlags) != 0 {
			key = *v.State + "+" + strings.Join(*v.StateFlags, "+")
		} else {
			key = *v.State
		}
		//t.Stats.StateCnt[*v.JobState]++
		t.Stats.StateCnt[key]++
		t.Stats.StateSimpleCnt[*v.State]++
	}
	l.Printf("GetStatsFiltered end\n")
}
