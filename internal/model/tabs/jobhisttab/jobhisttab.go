package jobhisttab

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander-dev/internal/stats"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

type JobHistTab struct {
	StatsOn           bool
	HistFetched       bool // signals View() if sacct call is finished, to print "waiting for..." message
	HistFetchFail     bool // if sacct call times out/errors, this is set to true
	SacctTable        table.Model
	SacctHist         SacctJSON
	SacctHistFiltered SacctJSON
	Filter            textinput.Model
	Stats
}

type Stats struct {
	StateCnt map[string]uint
	AvgWait  time.Duration
	MinWait  time.Duration
	MaxWait  time.Duration
	MedWait  time.Duration
	AvgRun   time.Duration
	MinRun   time.Duration
	MaxRun   time.Duration
	MedRun   time.Duration
	SDWait   int
}

func (t *JobHistTab) GetStatsFiltered(l *log.Logger) {
	t.Stats.StateCnt = map[string]uint{}
	tmp := []time.Duration{}
	tmpRun := []time.Duration{}
	t.AvgWait = 0
	t.MedWait = 0

	l.Printf("GetStatsFiltered start on %d rows\n", len(t.SacctHistFiltered.Jobs))

	for _, v := range t.SacctHistFiltered.Jobs {
		t.Stats.StateCnt[*v.State.Current]++
		//l.Printf("TIME: submit=%d, start=%d, end=%d\n", *v.Time.Submission, *v.Time.Start, *v.Time.End)
		if *v.State.Current != "RUNNING" {
			tmp = append(tmp, time.Unix(int64(*v.Time.Start), 0).Sub(time.Unix(int64(*v.Time.Submission), 0)))
			tmpRun = append(tmpRun, time.Unix(int64(*v.Time.End), 0).Sub(time.Unix(int64(*v.Time.Start), 0)))
		} else {
			// TODO: Question, do we include RUNNING jobs in the calculation?
			//tmpRun = append(tmpRun, time.Duration(*v.Time.Elapsed))
		}
	}

	t.MedWait, t.MinWait, t.MaxWait = stats.Median(tmp)
	t.MedRun, t.MinRun, t.MaxRun = stats.Median(tmpRun)
	t.AvgWait = stats.Avg(tmp)
	t.AvgRun = stats.Avg(tmpRun)

	l.Printf("GetStatsFiltered end\n")
}
