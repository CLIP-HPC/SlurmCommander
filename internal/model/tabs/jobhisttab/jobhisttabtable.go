package jobhisttab

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/CLIP-HPC/SlurmCommander/internal/command"
	"github.com/CLIP-HPC/SlurmCommander/internal/slurm"
	"github.com/CLIP-HPC/SlurmCommander/internal/table"
)

var SacctTabCols = []table.Column{
	{
		Title: "Job ID",
		Width: 10,
	},
	{
		Title: "Job Name",
		Width: 35,
	},
	{
		Title: "Part.",
		Width: 5,
	},
	{
		Title: "QoS",
		Width: 10,
	},
	{
		Title: "Account",
		Width: 10,
	},
	{
		Title: "User",
		Width: 15,
	},
	{
		Title: "Nodes",
		Width: 20,
	},
	{
		Title: "State",
		Width: 10,
	},
}

type SacctJSON slurm.SacctJSON
type TableRows []table.Row

func (saList *SacctJSON) FilterSacctTable(f string, l *log.Logger) (*TableRows, *SacctJSON, *command.ErrorMsg) {
	var (
		saTabRows         = TableRows{}
		sacctHistFiltered = SacctJSON{}
		errMsg            *command.ErrorMsg
		re                *regexp.Regexp
	)

	l.Printf("FilterSacctTable: rows %d", len(saList.Jobs))
	re, err := regexp.Compile(f)
	if err != nil {
		l.Printf("FAIL: compile regexp: %q with err: %s", f, err)
		f = ""
		re, _ = regexp.Compile(f)
		errMsg = &command.ErrorMsg{
			From:    "FilterSacctTable",
			ErrHelp: "Regular expression failed to compile, please correct it (turn on DEBUG to see details)",
			OrigErr: err,
		}
	}

	for _, v := range saList.Jobs {

		// https://github.com/CLIP-HPC/SlurmCommander/issues/9
		// Entry without JobId, log&discard
		switch {
		case v.JobId == nil:
			l.Printf("FilterSacctTable: Found job with no JobId field, skipping...\n")
			continue
		}

		line := strings.Join([]string{
			strconv.Itoa(*v.JobId),
			*v.Name,
			*v.Qos,
			*v.Account,
			*v.User,
			*v.State.Current,
		}, ".")

		if re.MatchString(line) {
			saTabRows = append(saTabRows, table.Row{strconv.Itoa(*v.JobId), *v.Name, *v.Partition, *v.Qos, *v.Account, *v.User, *v.Nodes, *v.State.Current})
			sacctHistFiltered.Jobs = append(sacctHistFiltered.Jobs, v)
		}
	}

	return &saTabRows, &sacctHistFiltered, errMsg
}
