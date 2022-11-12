package jobhisttab

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/pja237/slurmcommander-dev/internal/slurm"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

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

type SacctJSON slurm.SacctJSON
type TableRows []table.Row

func (saList *SacctJSON) FilterSacctTable(f string, l *log.Logger) (TableRows, SacctJSON) {
	var (
		saTabRows         = TableRows{}
		sacctHistFiltered = SacctJSON{}
	)

	l.Printf("FilterSacctTable: rows %d", len(saList.Jobs))
	re, err := regexp.Compile(f)
	if err != nil {
		l.Printf("FAIL: compile regexp: %q with err: %s", f, err)
		f = ""
	}

	for _, v := range saList.Jobs {

		line := strings.Join([]string{
			strconv.Itoa(*v.JobId),
			*v.Name,
			*v.Account,
			*v.User,
			*v.State.Current,
		}, ".")

		if re.MatchString(line) {
			saTabRows = append(saTabRows, table.Row{strconv.Itoa(*v.JobId), *v.Name, *v.Partition, *v.Account, *v.User, *v.State.Current})
			sacctHistFiltered.Jobs = append(sacctHistFiltered.Jobs, v)
		}
	}

	return saTabRows, sacctHistFiltered
}
