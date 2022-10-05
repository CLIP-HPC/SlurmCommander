package model

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander/internal/command"
)

func (m Model) Init() tea.Cmd {
	// use bubbletea.Batch(com1, com2, ...) here
	//return command.TimedGetSqueue()
	//return tea.Batch(command.TimedGetSqueue(), command.TimedGetSinfo(), command.TimedGetSacct())
	return tea.Batch(command.QuickGetSqueue(), command.QuickGetSinfo(), command.QuickGetSacct())
}
