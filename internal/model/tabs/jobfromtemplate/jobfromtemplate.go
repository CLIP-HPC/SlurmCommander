package jobfromtemplate

import (
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/keybindings"
	"github.com/pja237/slurmcommander-dev/internal/table"
)

type JobFromTemplateTab struct {
	TemplatesTable table.Model
	TemplatesList  TemplatesListRows
	TemplateEditor textarea.Model
	NewJobScript   string
	EditTemplate   bool
}

type EditTemplate bool

func EditorOn() tea.Cmd {
	return func() tea.Msg {
		return EditTemplate(true)
	}
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
	&keybindings.DefaultKeyMap.Stats:         false,
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
	&keybindings.DefaultKeyMap.Stats:         false,
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
		"/etc/slurmcommander-dev/templates",
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

func SaveToFile(name string, content string, l *log.Logger) (string, error) {
	var (
		ofName string
	)
	ofName = strings.TrimSuffix(name, ".sbatch")
	ofName += "-" + strconv.FormatInt(time.Now().Unix(), 10)
	ofName += ".sbatch"
	l.Printf("SaveToFile INFO: OutputFileName %s\n", ofName)
	if err := os.WriteFile(ofName, []byte(content), 0644); err != nil {
		l.Printf("SaveToFile ERROR: File %s: %s\n", ofName, err)
		return "", err
	}
	return ofName, nil
}
