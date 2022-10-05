package slurm

import (
	"github.com/charmbracelet/bubbles/table"
	"github.com/pja237/slurmcommander/internal/openapi"
)

type SinfoJSON struct {
	Nodes []openapi.V0039Node
}

type NodeJSON struct {
	Name         string
	Address      string
	Hostname     string
	State        string
	State_flags  []string
	Cpus         int
	Alloc_cpus   int
	Idle_cpus    int
	Alloc_memory int
	Real_memory  int
	Free_memory  int
}

var SinfoTabCols = []table.Column{
	{
		Title: "Name",
		Width: 15,
	},
	{
		Title: "State",
		Width: 10,
	},
	{
		Title: "Cpus",
		Width: 4,
	},
	{
		Title: "Cpus Available",
		Width: 4,
	},
	{
		Title: "Memory Total",
		Width: 10,
	},
	{
		Title: "Memory Available",
		Width: 10,
	},
	{
		Title: "State FLAGS",
		Width: 40,
	},
}

// TODO: not sure what i was thinking, but this we really don't need, just inject in Update() directly to model!?
var SinfoTabRows = []table.Row{}
