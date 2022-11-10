package slurm

import (
	"log"
	"regexp"
	"strconv"
	"time"

	"github.com/pja237/slurmcommander-dev/internal/openapi"
	"github.com/pja237/slurmcommander-dev/internal/table"
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
	{
		Title: "Priority",
		Width: 10,
	},
}

type TableRows []table.Row

func (sqJson *SqueueJSON) FilterSqueueTable(f string, l *log.Logger) (TableRows, SqueueJSON) {
	var (
		sqTabRows      = TableRows{}
		sqJsonFiltered = SqueueJSON{}
	)

	// TODO:  this is too slow, join & re2

	l.Printf("Filter SQUEUE start.\n")
	re, err := regexp.Compile(f)
	if err != nil {
		l.Printf("FAIL: compile regexp: %q with err: %s", f, err)
		f = ""
	}
	t := time.Now()
	for _, v := range sqJson.Jobs {
		app := false
		if f != "" {
			switch {
			case re.MatchString(strconv.Itoa(*v.JobId)):
				app = true
			case re.MatchString(*v.Name):
				app = true
			case re.MatchString(*v.Account):
				app = true
			case re.MatchString(*v.UserName):
				app = true
			case re.MatchString(*v.JobState):
				app = true
			}
		} else {
			app = true
		}
		if app {
			sqTabRows = append(sqTabRows, table.Row{strconv.Itoa(*v.JobId), *v.Name, *v.Account, *v.UserName, *v.JobState, strconv.Itoa(*v.Priority)})
			sqJsonFiltered.Jobs = append(sqJsonFiltered.Jobs, v)
		}
	}
	l.Printf("Filter SQUEUE end in %.3f seconds\n", time.Since(t).Seconds())

	return sqTabRows, sqJsonFiltered
}
