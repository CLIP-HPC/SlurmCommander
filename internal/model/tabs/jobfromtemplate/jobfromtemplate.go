package jobfromtemplate

import (
	"log"
	"os"

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
	&keybindings.DefaultKeyMap.Quit:          false,
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
	{
		Title: "Path",
		Width: 30,
	},
}

type TemplateText string

func GetTemplate(name string, l *log.Logger) tea.Cmd {
	l.Printf("GetTemplate: open file: %s\n", name)
	return func() tea.Msg {

		t, e := os.ReadFile(name)
		if e != nil {
			l.Printf("GetTemplate ERROR: open(%s): %s", name, e)

		}
		return TemplateText(t)
	}
}

type TemplatesListRows []table.Row

var (
	DefaultTemplatePaths = []string{
		"/etc/slurmcommander/templates",
	}
)

func GetTemplateList(paths []string, l *log.Logger) tea.Cmd {

	return func() tea.Msg {
		var tlr TemplatesListRows
		for _, p := range paths {
			files, err := os.ReadDir(p)
			if err != nil {
				l.Printf("GetTemplateList ERROR: %s\n", err)
			}
			for _, f := range files {
				l.Printf("GetTemplateList INFO files: %s %s\n", p, f.Name())
				// TODO: check suffix, if .sbatch, append to tlr
				// if suffix=".descr" then read content and use as description
				tlr = append(tlr, table.Row{f.Name(), "Description", p + "/" + f.Name()})
			}

		}
		return tlr
	}

}
