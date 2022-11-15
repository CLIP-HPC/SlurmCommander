package jobtab

import (
	"log"
	"sort"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander-dev/internal/stats"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

type JobTab struct {
	InfoOn           bool
	CountsOn         bool
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
	Top5user   CountItemSlice
	Top5acc    CountItemSlice
	JobPerQos  CountItemSlice
	JobPerPart CountItemSlice
}

type CountItemSlice []CountItem

type CountItem struct {
	Name  string
	Count uint
}

type CountItemMap map[string]uint

func sortItemMap(m *CountItemMap) CountItemSlice {
	var ret = CountItemSlice{}
	//ret := make(CountItemSlice, len(*m))
	for k, v := range *m {
		ret = append(ret, CountItem{
			Name:  k,
			Count: v,
		})
	}

	sort.Slice(ret, func(i, j int) bool {
		if ret[i].Count > ret[j].Count {
			return true
		} else {
			return false
		}
	})

	return ret
}

func top5(src CountItemSlice) CountItemSlice {
	var ret CountItemSlice
	for i, v := range src {
		if i < 5 {
			ret = append(ret, v)
		}
	}
	return ret
}

func (t *JobTab) GetStatsFiltered(l *log.Logger) {

	top5user := CountItemMap{}
	top5acc := CountItemMap{}
	jpq := CountItemMap{}
	jpp := CountItemMap{}

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
		top5acc[*v.Account]++
		top5user[*v.UserName]++
		jpp[*v.Partition]++
		jpq[*v.Qos]++
	}

	// sort & filter breakdowns
	t.Breakdowns.Top5user = top5(sortItemMap(&top5user))
	t.Breakdowns.Top5acc = top5(sortItemMap(&top5acc))
	t.Breakdowns.JobPerPart = sortItemMap(&jpp)
	t.Breakdowns.JobPerQos = sortItemMap(&jpq)

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
