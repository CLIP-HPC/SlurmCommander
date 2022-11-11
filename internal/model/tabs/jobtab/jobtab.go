package jobtab

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander-dev/internal/stats"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

type JobTab struct {
	InfoOn           bool
	StatsOn          bool
	SqueueTable      table.Model
	Squeue           SqueueJSON
	SqueueFiltered   SqueueJSON
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
			tmpRun = append(tmpRun, time.Since(time.Unix(int64(*v.StartTime), 0)))
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
