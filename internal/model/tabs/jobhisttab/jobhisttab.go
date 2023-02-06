package jobhisttab

import (
	"log"
	"time"
	"strconv"

	"github.com/charmbracelet/bubbles/textinput"
	"github.com/CLIP-HPC/SlurmCommander/internal/generic"
	"github.com/CLIP-HPC/SlurmCommander/internal/stats"
	"github.com/CLIP-HPC/SlurmCommander/internal/table"
)

type JobHistTab struct {
	StatsOn           bool
	CountsOn          bool
	FilterOn          bool
	UserInputsOn      bool // allow user to add/modify parameters for slurm commands
	HistFetched       bool // signals View() if sacct call is finished, to print "waiting for..." message
	HistFetchFail     bool // if sacct call times out/errors, this is set to true
	JobHistStart      uint
	JobHistTimeout    uint
	SacctTable        table.Model
	SacctHist         SacctJSON
	SacctHistFiltered SacctJSON
	Filter            textinput.Model
	UserInputs        generic.UserInputs
	SacctCmdline      string // with format specifiers
	Stats
	Breakdowns
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

type Breakdowns struct {
	Top5user   generic.CountItemSlice
	Top5acc    generic.CountItemSlice
	JobPerQos  generic.CountItemSlice
	JobPerPart generic.CountItemSlice
}

func NewUserInputs(d uint, t uint, cmdline string) generic.UserInputs {
	var tmp_s string
	var tmp_t textinput.Model

	// User Input (for modifying the table/view)
	userinput := generic.UserInputs {
		FocusIndex: 0,
		ParamTexts: make([]string, 3),
		Params: make([]textinput.Model, 3),
	}

	// TODO: we should probably think of encapsulating this
	for i := range userinput.Params {
		tmp_t = textinput.New()

		// FIXME we need to generalise this
		switch i {
		case 0:
		tmp_t.Placeholder = "Days"
		tmp_t.SetValue(strconv.FormatInt(int64(d), 10))
		tmp_t.Focus()
		tmp_t.CharLimit = 30
		tmp_t.Width = 30
		tmp_s = "Days"

		case 1:
		tmp_t.Placeholder = "Timeout (s)"
		tmp_t.SetValue(strconv.FormatInt(int64(t), 10))
		tmp_t.CharLimit = 30
		tmp_t.Width = 30
		tmp_s = "Timeout"

		case 2:
		tmp_t.Placeholder = "Cmdline"
		tmp_t.SetValue(cmdline)
		tmp_t.CharLimit = 50
		tmp_t.Width = 50
		tmp_s = "Cmdline"

		}

		userinput.Params[i] = tmp_t
		userinput.ParamTexts[i] = tmp_s
	}

	return userinput
}

func (t *JobHistTab) AdjTableHeight(h int, l *log.Logger) {
	l.Printf("FixTableHeight(%d) from %d\n", h, t.SacctTable.Height())
	if t.CountsOn || t.FilterOn || t.UserInputsOn {
		t.SacctTable.SetHeight(h - 30)
	} else {
		t.SacctTable.SetHeight(h - 15)
	}
	l.Printf("FixTableHeight to %d\n", t.SacctTable.Height())
}

func (t *JobHistTab) GetStatsFiltered(l *log.Logger) {
	top5user := generic.CountItemMap{}
	top5acc := generic.CountItemMap{}
	jpq := generic.CountItemMap{}
	jpp := generic.CountItemMap{}

	t.Stats.StateCnt = map[string]uint{}
	tmp := []time.Duration{}    // waiting times
	tmpRun := []time.Duration{} // running times
	t.AvgWait = 0
	t.MedWait = 0

	l.Printf("GetStatsFiltered start on %d rows\n", len(t.SacctHistFiltered.Jobs))

	for _, v := range t.SacctHistFiltered.Jobs {
		t.Stats.StateCnt[*v.State.Current]++
		//l.Printf("TIME: submit=%d, start=%d, end=%d\n", *v.Time.Submission, *v.Time.Start, *v.Time.End)
		switch *v.State.Current {
		case "PENDING":
			// no *v.Time.Start
			tmp = append(tmp, time.Since(time.Unix(int64(*v.Time.Submission), 0)))
		case "RUNNING":
			// no *v.Time.End
			tmpRun = append(tmpRun, time.Since(time.Unix(int64(*v.Time.Start), 0)))
		default:
			tmp = append(tmp, time.Unix(int64(*v.Time.Start), 0).Sub(time.Unix(int64(*v.Time.Submission), 0)))
			tmpRun = append(tmpRun, time.Unix(int64(*v.Time.End), 0).Sub(time.Unix(int64(*v.Time.Start), 0)))

		}
		// Breakdowns:
		if _, ok := top5acc[*v.Account]; !ok {
			top5acc[*v.Account] = &generic.CountItem{}
		}
		if _, ok := top5user[*v.User]; !ok {
			top5user[*v.User] = &generic.CountItem{}
		}
		if _, ok := jpp[*v.Partition]; !ok {
			jpp[*v.Partition] = &generic.CountItem{}
		}
		if _, ok := jpq[*v.Qos]; !ok {
			jpq[*v.Qos] = &generic.CountItem{}
		}
		top5acc[*v.Account].Count++
		top5user[*v.User].Count++
		jpp[*v.Partition].Count++
		jpq[*v.Qos].Count++
	}

	// sort & filter breakdowns
	t.Breakdowns.Top5user = generic.Top5(generic.SortItemMapBySel("Count", &top5user))
	t.Breakdowns.Top5acc = generic.Top5(generic.SortItemMapBySel("Count", &top5acc))
	t.Breakdowns.JobPerPart = generic.SortItemMapBySel("Count", &jpp)
	t.Breakdowns.JobPerQos = generic.SortItemMapBySel("Count", &jpq)

	t.MedWait, t.MinWait, t.MaxWait = stats.Median(tmp)
	t.MedRun, t.MinRun, t.MaxRun = stats.Median(tmpRun)
	t.AvgWait = stats.Avg(tmp)
	t.AvgRun = stats.Avg(tmpRun)

	l.Printf("GetStatsFiltered end\n")
}
