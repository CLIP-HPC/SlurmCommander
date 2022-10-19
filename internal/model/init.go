package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander/internal/command"
	"github.com/pja237/slurmcommander/internal/model/tabs/jobfromtemplate"
)

func (m Model) Init() tea.Cmd {
	// use bubbletea.Batch(com1, com2, ...) here
	//return command.TimedGetSqueue()
	//return tea.Batch(command.TimedGetSqueue(), command.TimedGetSinfo(), command.TimedGetSacct())
	return tea.Batch(
		command.GetUserName(m.Log),
		command.QuickGetSqueue(),
		command.QuickGetSinfo(),
		jobfromtemplate.GetTemplateList(jobfromtemplate.DefaultTemplatePaths, m.Log),
	)
}
