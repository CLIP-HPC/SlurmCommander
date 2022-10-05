package slurm

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/pja237/slurmcommander/internal/openapi"
)

type SqueueJSON struct {
	Jobs []openapi.V0039JobResponseProperties
}

type JobJSON struct {
	Account                   string
	Job_state                 string
	Job_id                    uint64
	User_name                 string
	Name                      string
	Command                   string
	Cpus                      int
	Node_count                int
	Tasks                     int
	Partition                 string
	Memory_per_node           int
	Priority                  int
	Resv_name                 string
	Qos                       string
	Time_limit                int
	Tres_req_str              string
	Current_working_directory string
}

// TODO: move table related data to model/tabs/jobtab.go
var SqueueTabCols = []table.Column{
	{
		Title: "Job ID",
		Width: 10,
	},
	{
		Title: "Job Name",
		Width: 20,
	},
	{
		Title: "Account",
		Width: 10,
	},
	{
		Title: "User Name",
		Width: 20,
	},
	{
		Title: "Job State",
		Width: 10,
	},
}

var SqueueTabRows = []table.Row{}
