package model

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/pja237/slurmcommander-dev/internal/generic"
	"github.com/pja237/slurmcommander-dev/internal/keybindings"
	"github.com/pja237/slurmcommander-dev/internal/styles"
	"github.com/pja237/slurmcommander-dev/internal/version"
)

// genTabs() generates top tabs
func (m Model) genTabs() string {

	var doc strings.Builder

	tlist := make([]string, len(tabs))
	for i, v := range tabs {
		if i == int(m.ActiveTab) {
			tlist = append(tlist, styles.TabActiveTab.Render(v))
		} else {
			tlist = append(tlist, styles.Tab.Render(v))
		}
	}
	row := lipgloss.JoinHorizontal(lipgloss.Top, tlist...)

	//gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	gap := styles.TabGap.Render(strings.Repeat(" ", max(0, m.winW-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	doc.WriteString(row + "\n")

	return doc.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) tabJobFromTemplate() string {

	if m.EditTemplate {
		return m.TemplateEditor.View()
	} else {
		if len(m.JobFromTemplateTab.TemplatesList) == 0 {
			return styles.NotFound.Render("\nNo templates found!\n")
		} else {
			return m.TemplatesTable.View()
		}
	}
}

func (m Model) tabAbout() string {

	s := "Version: " + version.BuildVersion + "\n"
	s += "Commit : " + version.BuildCommit + "\n"

	s += `
petar.jager@imba.oeaw.ac.at
CLIP-HPC Team @ VBC

Contributors:
`

	return s
}

func (m *Model) genTabHelp() string {
	var th string
	switch m.ActiveTab {
	case tabJobs:
		th = "List of jobs in the queue"
	case tabJobHist:
		th = fmt.Sprintf("List of jobs in the last %d days from all user associated accounts. (timeout: %d seconds)", m.JobHistTab.JobHistStart, m.JobHistTab.JobHistTimeout)
	case tabJobDetails:
		th = "Job details, select a job from Job History tab"
	case tabJobFromTemplate:
		th = "Edit and submit one of the job templates"
	case tabCluster:
		th = "List and status of cluster nodes"
	default:
		th = "SlurmCommander"
	}
	return th + "\n"
}

func (m Model) JobTabStats() string {

	m.Log.Printf("JobTabStats called\n")

	//str := "Queue statistics (filtered):\n\n"
	str := styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Job states (filtered):"))
	str += "\n\n"

	str += generic.GenCountStrVert(m.JobTab.Stats.StateCnt, m.Log)

	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Pending jobs:"))
	str += "\n\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinWait", generic.HumanizeDuration(m.JobTab.Stats.MinWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "AvgWait", generic.HumanizeDuration(m.JobTab.Stats.AvgWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MedWait", generic.HumanizeDuration(m.JobTab.Stats.MedWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MaxWait", generic.HumanizeDuration(m.JobTab.Stats.MaxWait, m.Log))

	str += "\n"
	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Running jobs:"))
	str += "\n\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinRun", generic.HumanizeDuration(m.JobTab.Stats.MinRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "AvgRun", generic.HumanizeDuration(m.JobTab.Stats.AvgRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MedRun", generic.HumanizeDuration(m.JobTab.Stats.MedRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MaxRun", generic.HumanizeDuration(m.JobTab.Stats.MaxRun, m.Log))

	return str
}

func (m Model) View() string {

	var (
		scr        strings.Builder
		MainWindow strings.Builder
	)

	// HEADER / TABS
	scr.WriteString(m.genTabs())
	scr.WriteString(m.genTabHelp())

	if m.Debug {
		// One debug line
		scr.WriteString(fmt.Sprintf("%s Width: %d Height: %d ErrorMsg: %s\n", styles.TextRed.Render("DEBUG ON:"), m.Globals.winW, m.Globals.winH, m.Globals.ErrorMsg))
	}

	if m.Globals.ErrorHelp != "" {
		scr.WriteString(styles.ErrorHelp.Render(fmt.Sprintf("ERROR: %s", m.Globals.ErrorHelp)))
		scr.WriteString("\n")
	}

	// PICK and RENDER ACTIVE TAB
	switch m.ActiveTab {
	case tabJobs:
		m.Log.Printf("CALL JobTab.View()\n")
		MainWindow.WriteString(m.JobTab.View(m.Log))

	case tabJobHist:
		m.Log.Printf("CALL JobHistTab.View()\n")
		MainWindow.WriteString(m.JobHistTab.View(m.Log))

	case tabJobDetails:
		m.Log.Printf("CALL JobDetailsTab.View()\n")
		MainWindow.WriteString(m.JobDetailsTab.View(&m.JobHistTab, m.Log))

	case tabJobFromTemplate:
		MainWindow.WriteString(m.tabJobFromTemplate())

	case tabCluster:
		m.Log.Printf("CALL ClusterTab.View()\n")
		MainWindow.WriteString(m.ClusterTab.View(m.Log))

	case tabAbout:
		MainWindow.WriteString(m.tabAbout())
	}

	// FOOTER
	scr.WriteString(lipgloss.JoinVertical(lipgloss.Left, styles.MainWindow.Render(MainWindow.String()), styles.HelpWindow.Render(m.Help.View(keybindings.DefaultKeyMap))))

	return scr.String()
}
