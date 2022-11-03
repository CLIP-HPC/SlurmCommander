package slurm

import (
	"log"
	"strconv"
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/pja237/slurmcommander-dev/internal/openapidb"
)

type SacctList []SacctEntry

//type SacctList [][]string

type SacctEntry []string

// SacctJobHist struct holds job history.
// Comes from unmarshalling sacct -A -S now-Xdays --json call.
type SacctJobHist struct {
	Jobs []openapidb.Dbv0037Job
}

// This is to distinguish in Update() the return from the jobhisttab and jobdetails tab
type SacctSingleJobHist struct {
	Jobs []openapidb.Dbv0037Job
}

type SacctJobJSON struct {
	Account   string
	Job_id    int
	Nodes     string
	Partition string
	Priority  int
	Qos       string
}

var SacctTabCols = []table.Column{
	{
		Title: "Job ID",
		Width: 10,
	},
	{
		Title: "Job Name",
		Width: 30,
	},
	{
		Title: "Partition",
		Width: 10,
	},
	{
		Title: "State",
		Width: 10,
	},
	{
		Title: "Exit Code",
		Width: 10,
	},
}

func (saList *SacctJobHist) FilterSacctTable(f string, l *log.Logger) (TableRows, SacctJobHist) {
	var (
		saTabRows         = TableRows{}
		sacctHistFiltered = SacctJobHist{}
	)

	l.Printf("FilterSacctTable: rows %d", len(saList.Jobs))
	for _, v := range saList.Jobs {
		app := false
		if f != "" {
			switch {
			case strings.Contains(strconv.Itoa(*v.JobId), f):
				// Id
				app = true
			case strings.Contains(*v.Name, f):
				// Name
				app = true
			case strings.Contains(*v.State.Current, f):
				// State
				app = true
			}
		} else {
			app = true
		}
		if app {
			//saTabRows = append(saTabRows, table.Row{v[0], v[1], v[2], v[3], v[4]})
			saTabRows = append(saTabRows, table.Row{strconv.Itoa(*v.JobId), *v.Name, *v.Partition, *v.State.Current, *v.ExitCode.Status})
			sacctHistFiltered.Jobs = append(sacctHistFiltered.Jobs, v)
		}
	}

	return saTabRows, sacctHistFiltered
}
