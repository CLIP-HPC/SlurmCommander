package jobfromtemplate

import (
	"log"
	"strings"

	"github.com/CLIP-HPC/SlurmCommander/internal/styles"
)

func (jft *JobFromTemplateTab) tabJobFromTemplate() string {

	if jft.EditTemplate {
		return jft.TemplateEditor.View()
	} else {
		if len(jft.TemplatesList) == 0 {
			return styles.NotFound.Render("\nNo templates found!\n")
		} else {
			return jft.TemplatesTable.View()
		}
	}
}

func (jft *JobFromTemplateTab) View(l *log.Logger) string {
	var (
		MainWindow strings.Builder
	)

	MainWindow.WriteString(jft.tabJobFromTemplate())

	return MainWindow.String()
}
