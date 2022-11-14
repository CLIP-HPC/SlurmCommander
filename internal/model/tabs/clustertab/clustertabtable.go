package clustertab

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/pja237/slurmcommander-dev/internal/command"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

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

type SinfoJSON slurm.SinfoJSON
type TableRows []table.Row

func (siJson *SinfoJSON) FilterSinfoTable(f string, l *log.Logger) (*TableRows, *SinfoJSON, *command.ErrorMsg) {
	var (
		siTabRows      = TableRows{}
		siJsonFiltered = SinfoJSON{}
		errMsg         *command.ErrorMsg
		re             *regexp.Regexp
	)

	l.Printf("FilterSinfoTable: rows %d", len(siJson.Nodes))
	re, err := regexp.Compile(f)
	if err != nil {
		l.Printf("FAIL: compile regexp: %q with err: %s", f, err)
		f = ""
		re, _ = regexp.Compile(f)
		errMsg = &command.ErrorMsg{
			From:    "FilterSinfoTable",
			ErrHelp: "Regular expression failed to compile, please correct it (turn on DEBUG to see details)",
			OrigErr: err,
		}
	}

	for _, v := range siJson.Nodes {

		line := strings.Join([]string{
			*v.Name,
			*v.State,
			strings.Join(*v.StateFlags, ","),
		}, ".")

		if re.MatchString(line) {
			siTabRows = append(siTabRows, table.Row{*v.Name, strings.Join(*v.Partitions, ","), *v.State, strconv.FormatInt(*v.IdleCpus, 10), strconv.Itoa(*v.Cpus), strconv.Itoa(*v.FreeMemory), strconv.Itoa(*v.RealMemory), strings.Join(*v.StateFlags, ",")})
			siJsonFiltered.Nodes = append(siJsonFiltered.Nodes, v)
		}
	}

	return &siTabRows, &siJsonFiltered, errMsg
}

// TODO: not sure what i was thinking, but this we really don't need, just inject in Update() directly to model!?
var SinfoTabRows = []table.Row{}
