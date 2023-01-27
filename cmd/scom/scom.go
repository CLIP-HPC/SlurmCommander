package main

import (
	"fmt"
	"log"
	"os"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/charmbracelet/lipgloss"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/CLIP-HPC/SlurmCommander/internal/cmdline"
	"github.com/CLIP-HPC/SlurmCommander/internal/command"
	"github.com/CLIP-HPC/SlurmCommander/internal/config"
	"github.com/CLIP-HPC/SlurmCommander/internal/logger"
	"github.com/CLIP-HPC/SlurmCommander/internal/model"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/clustertab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobdetailstab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobfromtemplate"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobhisttab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobtab"
	"github.com/CLIP-HPC/SlurmCommander/internal/styles"
	"github.com/CLIP-HPC/SlurmCommander/internal/table"
	"github.com/CLIP-HPC/SlurmCommander/internal/version"
)

func main() {

	var (
		debugSet bool = false
		args *cmdline.CmdArgs
	)

	fmt.Printf("Welcome to Slurm Commander!\n\n")

	cc := config.NewConfigContainer()
	err := cc.GetConfig()
	if err != nil {
		log.Printf("ERROR: parsing config files: %s\n", err)
	}

	args, err = cmdline.NewCmdArgs()
	if err != nil {
		log.Fatalf("ERROR: parsing cmdline args: %s\n", err)
	}

	if *args.Version {
		version.DumpVersion()
		os.Exit(0)
	}

	// TODO: JFT We have the CMDline switches and config, now overwrite/append what's changed
	//log.Println(cc.DumpConfig())

	log.Printf("INFO: %s\n", cc.DumpConfig())
	// TODO: this is ugly, but quick. Rework, use model...
	command.NewCmdCC(*cc)
	jobtab.NewCmdCC(*cc)
	clustertab.NewCmdCC(*cc)
	jobhisttab.NewCmdCC(*cc)

	// TODO: move all this away to view/styles somewhere...
	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(styles.Bluegrey).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Background(styles.Blue).
		Foreground(styles.Yellow).
		Bold(false)

	// Filter TextInput
	ti := textinput.New()
	ti.Placeholder = ""
	ti.Focus()
	ti.CharLimit = 30
	ti.Width = 30

	// logging
	debugSet, l := logger.SetupLogger()

	// setup help
	hlp := help.New()
	hlp.Styles.ShortKey = styles.TextYellow
	hlp.Styles.ShortDesc = styles.TextBlue
	hlp.Styles.ShortSeparator = styles.TextBlueGrey

	/// JD viewport
	vp := viewport.New(10, 10)
	vp.Style = styles.JDviewportBox

	m := model.Model{
		Globals: model.Globals{
			Help:            hlp,
			ActiveTab:       0,
			Log:             l,
			Debug:           debugSet,
			ConfigContainer: *cc,
		},
		JobTab: jobtab.JobTab{
			SqueueTable: table.New(table.WithColumns(jobtab.SqueueTabCols), table.WithRows(jobtab.TableRows{}), table.WithStyles(s)),
			Filter: ti,
		},
		JobHistTab: jobhisttab.JobHistTab{
			SacctTable:     table.New(table.WithColumns(jobhisttab.SacctTabCols), table.WithRows(jobtab.TableRows{}), table.WithStyles(s)),
			Filter: ti,
			UserInputs:     jobhisttab.NewUserInputs(cc.HistDays, cc.HistTimeout),
			HistFetched:    false,
			HistFetchFail:  false,
			JobHistStart:   cc.HistDays,
			JobHistTimeout: cc.HistTimeout,
		},
		JobDetailsTab: jobdetailstab.JobDetailsTab{
			SelJobIDNew: -1,
			ViewPort:    vp,
		},
		JobFromTemplateTab: jobfromtemplate.JobFromTemplateTab{
			EditTemplate: false,
			NewJobScript: "",
			TemplatesTable: table.New(
				table.WithColumns(jobfromtemplate.TemplatesListCols),
				table.WithRows(jobfromtemplate.TemplatesListRows{}),
				table.WithStyles(s),
			),
		},
		ClusterTab: clustertab.ClusterTab{
			SinfoTable: table.New(table.WithColumns(clustertab.SinfoTabCols), table.WithRows(jobtab.TableRows{}), table.WithStyles(s)),
			Filter: ti,
		},
	}

	// OLD: bubbletea@v0.22.1
	//p := tea.NewProgram(tea.Model(m), tea.WithAltScreen())
	//if err := p.Start(); err != nil {
	//	log.Fatalf("ERROR: starting tea program: %q\n", err)
	//}

	// NEW: bubbletea@v0.23.0
	// Run returns the model as a tea.Model.
	p := tea.NewProgram(tea.Model(m), tea.WithAltScreen())
	ret, err := p.Run()
	if err != nil {
		fmt.Printf("Error starting program: %s", err)
		os.Exit(1)
	}

	if retMod, ok := ret.(model.Model); ok && retMod.Globals.SizeErr != "" {
		fmt.Printf("%s\n", retMod.Globals.SizeErr)
	}
	fmt.Printf("Goodbye!\n")
}
