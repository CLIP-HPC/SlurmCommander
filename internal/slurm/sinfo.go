package slurm

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/pja237/slurmcommander-dev/internal/openapi"
	"github.com/pja237/slurmcommander-dev/internal/table"
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
		Title: "Part.",
		Width: 5,
	},
	{
		Title: "State",
		Width: 10,
	},
	{
		Title: "CPUAvail",
		Width: 10,
	},
	{
		Title: "CPUTotal",
		Width: 10,
	},
	{
		Title: "MEMAvail",
		Width: 10,
	},
	{
		Title: "MEMTotal",
		Width: 10,
	},
	{
		Title: "State FLAGS",
		Width: 40,
	},
}

func (siJson *SinfoJSON) FilterSinfoTable(f string, l *log.Logger) (TableRows, SinfoJSON) {
	var (
		siTabRows      = TableRows{}
		siJsonFiltered = SinfoJSON{}
	)

	l.Printf("FilterSinfoTable: rows %d", len(siJson.Nodes))
	re, err := regexp.Compile(f)
	if err != nil {
		l.Printf("FAIL: compile regexp: %q with err: %s", f, err)
		f = ""
	}
	for _, v := range siJson.Nodes {
		app := false
		if f != "" {
			switch {
			case re.MatchString(*v.Name):
				app = true
			case re.MatchString(*v.State):
				app = true
			case re.MatchString(strings.Join(*v.StateFlags, ",")):
				app = true
			}
		} else {
			app = true
		}
		if app {
			siTabRows = append(siTabRows, table.Row{*v.Name, strings.Join(*v.Partitions, ","), *v.State, strconv.FormatInt(*v.IdleCpus, 10), strconv.Itoa(*v.Cpus), strconv.Itoa(*v.FreeMemory), strconv.Itoa(*v.RealMemory), strings.Join(*v.StateFlags, ",")})
			siJsonFiltered.Nodes = append(siJsonFiltered.Nodes, v)
		}
	}

	return siTabRows, siJsonFiltered
}

// TODO: not sure what i was thinking, but this we really don't need, just inject in Update() directly to model!?
var SinfoTabRows = []table.Row{}
