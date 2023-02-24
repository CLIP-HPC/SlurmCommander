package clustertab

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/CLIP-HPC/SlurmCommander/internal/command"
	"github.com/CLIP-HPC/SlurmCommander/internal/slurm"
	"github.com/CLIP-HPC/SlurmCommander/internal/styles"
	"github.com/CLIP-HPC/SlurmCommander/internal/table"
	"github.com/charmbracelet/lipgloss"
)

const (
	// Width of the SinfoTable, used in calculating Stats box width.
	// Must be adjusted alongside SinfoTabCols changes.
	//SinfoTabWidth = 118
	SinfoTabWidth = 134
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
		Width: 8,
	},
	{
		Title: "CPUTotal",
		Width: 8,
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
		Title: "GPUAvail",
		Width: 8,
	},
	{
		Title: "GPUTotal",
		Width: 8,
	},
	{
		Title: "State FLAGS",
		Width: 15,
	},
}

type SinfoJSON slurm.SinfoJSON
type TableRows []table.Row

type CellStyle struct{}
type RowStyle struct{}
type SelectedStyle struct{}

func (h CellStyle) Style(i int, r table.Row) lipgloss.Style {
	switch r[i] {
	case "down":
		return styles.TextRed
	case "mixed":
		return styles.TextOrange
	default:
		return lipgloss.NewStyle()
	}
}

func (h RowStyle) Style(i int, r table.Row) lipgloss.Style {
	return lipgloss.NewStyle()
}

func (h SelectedStyle) Style(i int, r table.Row) lipgloss.Style {
	return styles.SelectedRow
}

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
			strings.Join(*v.Partitions, ","),
			*v.State,
			strings.Join(*v.StateFlags, ","),
		}, ".")

		if re.MatchString(line) {
			// This is how many GPUs are available on the node
			gpuAvail := slurm.ParseGRES(*v.Gres)

			// This is how many GPUs are allocated on the node
			gpuAlloc := slurm.ParseGRES(*v.GresUsed)
			siTabRows = append(siTabRows, table.Row{*v.Name, strings.Join(*v.Partitions, ","), *v.State, strconv.FormatInt(*v.IdleCpus, 10), strconv.Itoa(*v.Cpus), strconv.Itoa(*v.FreeMemory), strconv.Itoa(*v.RealMemory), strconv.Itoa(*gpuAvail - *gpuAlloc), strconv.Itoa(*gpuAvail), strings.Join(*v.StateFlags, ",")})
			siJsonFiltered.Nodes = append(siJsonFiltered.Nodes, v)
		}
	}

	return &siTabRows, &siJsonFiltered, errMsg
}

// TODO: not sure what i was thinking, but this we really don't need, just inject in Update() directly to model!?
var SinfoTabRows = []table.Row{}
