package slurm

import (
	"strings"

	"github.com/charmbracelet/bubbles/table"
	"github.com/pja237/slurmcommander/internal/openapidb"
)

type SacctList []SacctEntry

//type SacctList [][]string

type SacctEntry []string

// SacctJob struct describes specific job.
// Comes frofm unmarshalling sacct -j X --json call.
type SacctJob struct {
	//Jobs []openapidb.Dbv0039Job
	Jobs []openapidb.Dbv0037Job
	// TODO: incorporate dbv0.0.39 where sacct response is defined.
	// This is not the response from sacct.
	//Jobs []openapi.V0039JobResponseProperties
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
		Width: 15,
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
		Width: 20,
	},
	{
		Title: "Exit Code",
		Width: 10,
	},
}

func (saList *SacctList) FilterSacctTable(f string) (TableRows, SacctList) {
	var (
		saTabRows      = TableRows{}
		saListFiltered = SacctList{}
	)

	for _, v := range *saList {
		app := false
		if f != "" {
			switch {
			case strings.Contains(v[0], f):
				// Id
				app = true
			case strings.Contains(v[1], f):
				// Name
				app = true
			case strings.Contains(v[3], f):
				// State
				app = true
			}
		} else {
			app = true
		}
		if app {
			saTabRows = append(saTabRows, table.Row{v[0], v[1], v[2], v[3], v[4]})
			saListFiltered = append(saListFiltered, v)
		}
	}

	return saTabRows, saListFiltered
}
