package jobfromtemplate

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/defaults"
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

var TemplatesListCols = []table.Column{
	{
		Title: "Name",
		Width: 20,
	},
	{
		Title: "Description",
		Width: 40,
	},
	{
		Title: "Path",
		Width: 100,
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

func GetTemplateList(paths []string, l *log.Logger) tea.Cmd {

	return func() tea.Msg {
		var tlr TemplatesListRows
		for _, p := range paths {
			l.Printf("GetTemplateList reading dir: %s\n", p)
			files, err := os.ReadDir(p)
			if err != nil {
				l.Printf("GetTemplateList ERROR: %s\n", err)
			}
			for _, f := range files {
				l.Printf("GetTemplateList INFO files: %s %s\n", p, f.Name())
				// if suffix=".desc" then read content and use as description
				if strings.HasSuffix(f.Name(), ".sbatch") {
					sbatchPath := p + "/" + f.Name()
					descPath := p + "/" + strings.TrimSuffix(f.Name(), defaults.TemplatesSuffix) + defaults.TemplatesDescSuffix
					fd, err := os.Open(descPath)
					if err != nil {
						// handle error and put no desc in table
						l.Printf("GetTemplateList FAIL open desc file: %s\n", err)
						tlr = append(tlr, table.Row{f.Name(), "", sbatchPath})
					} else {
						l.Printf("GetTemplateList INFO open desc file: %s\n", descPath)
						s := bufio.NewScanner(fd)
						if s.Scan() {
							tlr = append(tlr, table.Row{f.Name(), s.Text(), sbatchPath})
						} else {
							l.Printf("GetTemplateList ERR scanning desc file: %s : %s\n", descPath, s.Err())
							// then we put no description
							tlr = append(tlr, table.Row{f.Name(), "", sbatchPath})
						}
					}
				}
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
