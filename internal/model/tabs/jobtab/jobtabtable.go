package jobtab

import (
	"log"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

var SqueueTabCols = []table.Column{
	{
		Title: "Job ID",
		Width: 10,
	},
	{
		Title: "Job Name",
		Width: 60,
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

type SqueueJSON slurm.SqueueJSON
type TableRows []table.Row

func (sqJson *SqueueJSON) FilterSqueueTable(f string, l *log.Logger) (*TableRows, *SqueueJSON, *command.ErrorMsg) {
	var (
		sqTabRows      = TableRows{}
		sqJsonFiltered = SqueueJSON{}
		errMsg         *command.ErrorMsg
		re             *regexp.Regexp
	)

	l.Printf("Filter SQUEUE start.\n")
	re, err := regexp.Compile(f)
	if err != nil {
		// User entered bad regexp, return error and set empty re (filter.value() will be zeroed in update())
		l.Printf("FAIL: compile regexp: %q with err: %s", f, err)
		f = ""
		re, _ = regexp.Compile(f)
		errMsg = &command.ErrorMsg{
			From:    "FilterSqueueTable",
			ErrHelp: "Regular expression failed to compile, please correct it (turn on DEBUG to see details)",
			OrigErr: err,
		}
	}
	t := time.Now()

	// NEW: Filter:
	// 	1. join strings from job into one line
	// 	2. re.MatchString(line)
	// Allows doing matches across multiple columns, e.g.: "bash.*(als|gmi)" - "bash jobs BY als or gmi accounts"
	l.Printf("Filtering on len=%d\n", len(sqJson.Jobs))
	//l.Printf("Filtering %#v\n", sqJson.Jobs)
	prio := 0
	for i, v := range sqJson.Jobs {
		l.Printf("Filtering %%%d\n", i)
		l.Printf("0. Filter SQTable: jobid %s\n", strconv.Itoa(*v.JobId))
		l.Printf("0. Filter SQTable: name %s\n", *v.Name)
		l.Printf("0. Filter SQTable: acc %s\n", *v.Account)
		l.Printf("0. Filter SQTable: username %s\n", *v.UserName)
		l.Printf("0. Filter SQTable: jobstate %s\n", *v.JobState)
		// https://groups.google.com/g/golang-nuts/c/i8JVM25NkDo
		if v.Priority == nil {
			l.Println("NIL PRIO")
			prio = 0
		} else {
			prio = *v.Priority
		}
		l.Printf("0. Filter SQTable: prio %s\n", strconv.Itoa(prio))

		// NEW: 1. JOIN
		line := strings.Join([]string{
			strconv.Itoa(*v.JobId),
			*v.Name,
			*v.Account,
			*v.UserName,
			*v.JobState,
		}, ".")

		// NEW: 2. MATCH
		if re.MatchString(line) {
			sqTabRows = append(sqTabRows, table.Row{strconv.Itoa(*v.JobId), *v.Name, *v.Account, *v.UserName, *v.JobState, strconv.Itoa(prio)})
			sqJsonFiltered.Jobs = append(sqJsonFiltered.Jobs, v)
		}
	}
	l.Printf("Filter SQUEUE end in %.3f seconds\n", time.Since(t).Seconds())

	return &sqTabRows, &sqJsonFiltered, errMsg
}
