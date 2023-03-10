package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/CLIP-HPC/SlurmCommander/internal/command"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/clustertab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobfromtemplate"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobtab"
)

func (m Model) Init() tea.Cmd {
	// use bubbletea.Batch(com1, com2, ...) here
	//return command.TimedGetSqueue()
	//return tea.Batch(command.TimedGetSqueue(), command.TimedGetSinfo(), command.TimedGetSacct())
	return tea.Batch(
		command.GetUserName(m.Log),
		jobtab.QuickGetSqueue(m.Log),
		clustertab.QuickGetSinfo(m.Log),
		jobfromtemplate.GetTemplateList(m.Globals.ConfigContainer.TemplateDirs, m.Log),
	)
}
