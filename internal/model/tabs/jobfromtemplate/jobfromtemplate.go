package jobfromtemplate

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander/internal/keybindings"
)

type JobFromTemplateTab struct {
	TemplatesTable table.Model
	TemplateEditor textarea.Model
	EditTemplate   bool
	NewJobScript   string
}

type Keys map[*key.Binding]bool

var oKeyMap = Keys{
	&keybindings.DefaultKeyMap.Up:       true,
	&keybindings.DefaultKeyMap.Down:     true,
	&keybindings.DefaultKeyMap.PageUp:   false,
	&keybindings.DefaultKeyMap.PageDown: false,
	&keybindings.DefaultKeyMap.Slash:    false,
	&keybindings.DefaultKeyMap.Info:     false,
	&keybindings.DefaultKeyMap.Enter:    true,
}

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.TtabSel:       true,
	&keybindings.DefaultKeyMap.Up:            true,
	&keybindings.DefaultKeyMap.Down:          true,
	&keybindings.DefaultKeyMap.PageUp:        false,
	&keybindings.DefaultKeyMap.PageDown:      false,
	&keybindings.DefaultKeyMap.Tab:           true,
	&keybindings.DefaultKeyMap.Slash:         false,
	&keybindings.DefaultKeyMap.Info:          false,
	&keybindings.DefaultKeyMap.Enter:         true,
	&keybindings.DefaultKeyMap.Quit:          true,
	&keybindings.DefaultKeyMap.SaveSubmitJob: false,
	&keybindings.DefaultKeyMap.Escape:        false,
}

var EditorKeyMap = Keys{
	&keybindings.DefaultKeyMap.TtabSel:       false,
	&keybindings.DefaultKeyMap.Up:            false,
	&keybindings.DefaultKeyMap.Down:          false,
	&keybindings.DefaultKeyMap.PageUp:        false,
	&keybindings.DefaultKeyMap.PageDown:      false,
	&keybindings.DefaultKeyMap.Tab:           false,
	&keybindings.DefaultKeyMap.Slash:         false,
	&keybindings.DefaultKeyMap.Info:          false,
	&keybindings.DefaultKeyMap.Enter:         false,
	&keybindings.DefaultKeyMap.Quit:          true,
	&keybindings.DefaultKeyMap.SaveSubmitJob: true,
	&keybindings.DefaultKeyMap.Escape:        true,
}

func (k *Keys) SetupKeys() {
	//for k, v := range KeyMap {
	for k, v := range *k {
		k.SetEnabled(v)
	}
}

var TemplatesListCols = []table.Column{
	{
		Title: "Name",
		Width: 20,
	},
	{
		Title: "Description",
		Width: 30,
	},
}

// TODO: this comes from template files
var TemplatesListRowsDef = []table.Row{
	{"Simple job", "Single node multiple threads"},
	{"Array job", "Multiple simple jobs"},
	{"Relion job v0.0.1", "Relion job with some parameters"},
	{"Relion job v0.0.2", "Relion job with some other parameters"},
}

type TemplateText string

func GetTemplate(name string) tea.Cmd {
	return func() tea.Msg {
		return TemplateSample
	}
}

var TemplateSample = TemplateText(`
#!/usr/bin/bash

#SBATCH --job-name=SimpleMultithreadJob
#SBATCH --ntasks=1 
#SBATCH --cpus-per-task=4  # 4 threads per task
#SBATCH --time=02:00:00    # two hours
#SBATCH --mem=1G           # 1 GB RAM per node
#SBATCH --output=.out

hostname
`)
