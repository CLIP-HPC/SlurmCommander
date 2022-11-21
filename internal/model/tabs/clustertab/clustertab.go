package clustertab

import (
	"log"
	"strings"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander-dev/internal/generic"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

type ClusterTab struct {
	StatsOn       bool
	CountsOn      bool
	FilterOn      bool
	SinfoTable    table.Model
	CpuBar        progress.Model
	MemBar        progress.Model
	GpuBar        progress.Model
	Sinfo         SinfoJSON
	SinfoFiltered SinfoJSON
	Filter        textinput.Model
	Stats
	Breakdowns
}

type Stats struct {
	// TODO: also perhaps: count by user? account?
	StateCnt       map[string]uint
	StateSimpleCnt map[string]uint
}

type Breakdowns struct {
	CpuPerPart    generic.CountItemSlice
	MemPerPart    generic.CountItemSlice
	GpuPerPart    generic.CountItemSlice
	NodesPerState generic.CountItemSlice
}

func (t *ClusterTab) AdjTableHeight(h int, l *log.Logger) {
	l.Printf("FixTableHeight(%d) from %d\n", h, t.SinfoTable.Height())
	if t.CountsOn || t.FilterOn {
		t.SinfoTable.SetHeight(h - 35)
	} else {
		t.SinfoTable.SetHeight(h - 25)
	}
	l.Printf("FixTableHeight to %d\n", t.SinfoTable.Height())
}

func (t *ClusterTab) GetStatsFiltered(l *log.Logger) {
	var key string

	cpp := generic.CountItemMap{} // CpuPerPartition
	mpp := generic.CountItemMap{} // MemPerPartition
	gpp := generic.CountItemMap{} // GPUPerPartition
	nps := generic.CountItemMap{} // NodesPerState

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

		//Breakdowns:
		// CpusPer
		for _, p := range *v.Partitions {
			if _, ok := cpp[p]; !ok {
				cpp[p] = &generic.CountItem{}
			}
			if _, ok := mpp[p]; !ok {
				mpp[p] = &generic.CountItem{}
			}
			if _, ok := gpp[p]; !ok {
				gpp[p] = &generic.CountItem{}
			}
			cpp[p].Name = p
			cpp[p].Count += uint(*v.AllocCpus)
			cpp[p].Total += uint(*v.Cpus)
			mpp[p].Name = p
			mpp[p].Count += uint(*v.AllocMemory)
			mpp[p].Total += uint(*v.RealMemory)
			gpp[p].Name = p
			gpp[p].Count += uint(*slurm.ParseGRES(*v.GresUsed))
			gpp[p].Total += uint(*slurm.ParseGRES(*v.Gres))
		}
		for _, s := range *v.StateFlags {
			if _, ok := nps[s]; !ok {
				nps[s] = &generic.CountItem{}
			}
			nps[s].Name = s
			nps[s].Count++
		}
	}

	// sort & filter breakdowns
	t.Breakdowns.CpuPerPart = generic.SortItemMapBySel("Name", &cpp)
	t.Breakdowns.MemPerPart = generic.SortItemMapBySel("Name", &mpp)
	t.Breakdowns.NodesPerState = generic.SortItemMapBySel("Count", &nps)
	t.Breakdowns.GpuPerPart = generic.SortItemMapBySel("Count", &gpp)

	l.Printf("GetStatsFiltered end\n")
}
