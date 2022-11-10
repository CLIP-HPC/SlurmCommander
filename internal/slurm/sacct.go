package slurm

import (
	"log"
	"regexp"
	"strconv"

	"github.com/pja237/slurmcommander-dev/internal/openapidb"
	"github.com/pja237/slurmcommander-dev/internal/table"
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
		Title: "Account",
		Width: 10,
	},
	{
		Title: "User",
		Width: 10,
	},
	{
		Title: "State",
		Width: 10,
	},
}

func (saList *SacctJobHist) FilterSacctTable(f string, l *log.Logger) (TableRows, SacctJobHist) {
	var (
		saTabRows         = TableRows{}
		sacctHistFiltered = SacctJobHist{}
	)

	l.Printf("FilterSacctTable: rows %d", len(saList.Jobs))
	re, err := regexp.Compile(f)
	if err != nil {
		l.Printf("FAIL: compile regexp: %q with err: %s", f, err)
		f = ""
	}
	for _, v := range saList.Jobs {
		app := false
		if f != "" {
			switch {
			case re.MatchString(strconv.Itoa(*v.JobId)):
				// Id
				app = true
			case re.MatchString(*v.Name):
				// Name
				app = true
			case re.MatchString(*v.Account):
				// State
				app = true
			case re.MatchString(*v.User):
				// State
				app = true
			case re.MatchString(*v.State.Current):
				// State
				app = true
			}
		} else {
			app = true
		}
		if app {
			//saTabRows = append(saTabRows, table.Row{v[0], v[1], v[2], v[3], v[4]})
			saTabRows = append(saTabRows, table.Row{strconv.Itoa(*v.JobId), *v.Name, *v.Partition, *v.Account, *v.User, *v.State.Current})
			sacctHistFiltered.Jobs = append(sacctHistFiltered.Jobs, v)
		}
	}

	return saTabRows, sacctHistFiltered
}
