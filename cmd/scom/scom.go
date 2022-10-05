package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pja237/slurmcommander/internal/model"
	"github.com/pja237/slurmcommander/internal/model/tabs/clustertab"
	"github.com/pja237/slurmcommander/internal/model/tabs/jobfromtemplate"
	"github.com/pja237/slurmcommander/internal/model/tabs/jobhisttab"
	"github.com/pja237/slurmcommander/internal/model/tabs/jobtab"
	"github.com/pja237/slurmcommander/internal/slurm"
)

func main() {

	var (
		logf *os.File
		err  error
	)

	fmt.Println("Welcome to Slurm Commander!")

	// TODO: move all this away to view somewhere...
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(lipgloss.Color("#0057B8")).
		Foreground(lipgloss.Color("#FFD700")).
		Bold(false)
	//s.Selected = s.Selected.
	//	Background(lipgloss.Color("#277BC0")).
	//	Foreground(lipgloss.Color("#FFCB42")).
	//	Bold(false)
	//
	//s.Selected = s.Selected.
	//	Foreground(lipgloss.Color("229")).
	//	Background(lipgloss.Color("57")).
	//	Bold(false)
	// Filter TextInput
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 20
	ti.Width = 20

	// logging
	if len(os.Getenv("DEBUG")) > 0 {
		logf, err = tea.LogToFile("debug.log", "debug")
		if err != nil {
			fmt.Println("fatal:", err)
			os.Exit(1)
		}
		logf.WriteString("Log to file set up.\n")
		//defer logf.Close()
	} else {
		os.OpenFile("/dev/null", os.O_WRONLY|os.O_APPEND, 0000)
	}

	m := model.Model{
		Globals: model.Globals{
			Help:         help.New(),
			ActiveTab:    0,
			LogF:         logf,
			FilterSwitch: -1,
		},
		JobTab: jobtab.JobTab{
			SqueueTable: table.New(table.WithColumns(slurm.SqueueTabCols), table.WithRows(slurm.SqueueTabRows), table.WithStyles(s)),
			Filter:      ti,
		},
		JobHistTab: jobhisttab.JobHistTab{
			SacctTable: table.New(table.WithColumns(slurm.SacctTabCols), table.WithRows(slurm.SacctTabRows), table.WithStyles(s)),
		},
		JobFromTemplateTab: jobfromtemplate.JobFromTemplateTab{
			EditTemplate: false,
			NewJobScript: "",
			TemplatesTable: table.New(
				table.WithColumns(jobfromtemplate.TemplatesListCols),
				table.WithRows(jobfromtemplate.TemplatesListRowsDef),
				table.WithStyles(s),
			),
		},
		JobClusterTab: clustertab.JobClusterTab{
			SinfoTable: table.New(table.WithColumns(slurm.SinfoTabCols), table.WithRows(slurm.SinfoTabRows), table.WithStyles(s)),
		},
	}

	logf.WriteString("Starting program.\n")
	//m.SqTable.SetStyles(s)
	p := tea.NewProgram(tea.Model(m), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatalf("ERROR: starting tea program: %q\n", err)
	}
	logf.WriteString("Good bye!\n")
}
