package slurm

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/pja237/slurmcommander/internal/openapidb"
)

type SacctList []SacctEntry

//type SacctList [][]string

type SacctEntry []string

// SacctJob struct describes specific job.
// Comes frofm unmarshalling sacct -j X --json call.
type SacctJob struct {
	Jobs []openapidb.Dbv0039Job
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

var SacctTabRows = []table.Row{}
