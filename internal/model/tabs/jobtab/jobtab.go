package jobtab

import (
	"log"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/CLIP-HPC/SlurmCommander/internal/generic"
	"github.com/CLIP-HPC/SlurmCommander/internal/stats"
	"github.com/CLIP-HPC/SlurmCommander/internal/table"
)

type JobTab struct {
	InfoOn           bool
	CountsOn         bool
	StatsOn          bool
	FilterOn         bool
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
	Breakdowns
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

type Breakdowns struct {
	Top5user   generic.CountItemSlice
	Top5acc    generic.CountItemSlice
	JobPerQos  generic.CountItemSlice
	JobPerPart generic.CountItemSlice
}

func (t *JobTab) AdjTableHeight(h int, l *log.Logger) {
	l.Printf("FixTableHeight(%d) from %d\n", h, t.SqueueTable.Height())
	if t.InfoOn || t.CountsOn || t.FilterOn {
		t.SqueueTable.SetHeight(h - 30)
	} else {
		t.SqueueTable.SetHeight(h - 15)
	}
	l.Printf("FixTableHeight to %d\n", t.SqueueTable.Height())
}

func (t *JobTab) GetStatsFiltered(l *log.Logger) {

	top5user := generic.CountItemMap{}
	top5acc := generic.CountItemMap{}
	jpq := generic.CountItemMap{}
	jpp := generic.CountItemMap{}

	t.Stats.StateCnt = map[string]uint{}
	tmp := []time.Duration{}
	tmpRun := []time.Duration{}
	t.AvgWait = 0
	t.MedWait = 0

	l.Printf("GetStatsFiltered start on %d rows\n", len(t.SqueueFiltered.Jobs))
	for _, v := range t.SqueueFiltered.Jobs {
		t.Stats.StateCnt[*v.JobState]++
		switch *v.JobState {
		case "PENDING":
			tmp = append(tmp, time.Since(time.Unix(int64(*v.SubmitTime), 0)))
		case "RUNNING":
			tmpRun = append(tmpRun, time.Since(time.Unix(int64(*v.StartTime), 0)))
		}

		// Breakdowns:
		if _, ok := top5acc[*v.Account]; !ok {
			top5acc[*v.Account] = &generic.CountItem{}
		}
		if _, ok := top5user[*v.UserName]; !ok {
			top5user[*v.UserName] = &generic.CountItem{}
		}
		if _, ok := jpp[*v.Partition]; !ok {
			jpp[*v.Partition] = &generic.CountItem{}
		}
		if _, ok := jpq[*v.Qos]; !ok {
			jpq[*v.Qos] = &generic.CountItem{}
		}
		top5acc[*v.Account].Count++
		top5user[*v.UserName].Count++
		jpp[*v.Partition].Count++
		jpq[*v.Qos].Count++
	}

	// sort & filter breakdowns
	t.Breakdowns.Top5user = generic.Top5(generic.SortItemMapBySel("Count", &top5user))
	t.Breakdowns.Top5acc = generic.Top5(generic.SortItemMapBySel("Count", &top5acc))
	t.Breakdowns.JobPerPart = generic.SortItemMapBySel("Count", &jpp)
	t.Breakdowns.JobPerQos = generic.SortItemMapBySel("Count", &jpq)

	//l.Printf("TOP5USER: %#v\n", t.Breakdowns.Top5user)
	//l.Printf("TOP5ACC: %#v\n", t.Breakdowns.Top5acc)
	//l.Printf("JobPerQos: %#v\n", t.Breakdowns.JobPerQos)
	//l.Printf("JobPerPart: %#v\n", t.Breakdowns.JobPerPart)

	t.MedWait, t.MinWait, t.MaxWait = stats.Median(tmp)
	t.MedRun, t.MinRun, t.MaxRun = stats.Median(tmpRun)
	t.AvgWait = stats.Avg(tmp)
	t.AvgRun = stats.Avg(tmpRun)

	l.Printf("GetStatsFiltered end\n")
}
