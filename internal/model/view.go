package model

import (
	"fmt"
	"strings"

	"github.com/CLIP-HPC/SlurmCommander/internal/keybindings"
	"github.com/CLIP-HPC/SlurmCommander/internal/styles"
	"github.com/CLIP-HPC/SlurmCommander/internal/version"
	"github.com/charmbracelet/lipgloss"
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

func (m Model) tabAbout() string {

	s := "Version: " + version.BuildVersion + "\n"
	s += "Commit : " + version.BuildCommit + "\n"

	s += `
Petar Jager
CLIP-HPC @VBC

Contributors:
Seren Ümit
Kilian Cavalotti
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

func (m Model) View() string {

	var (
		header     strings.Builder
		MainWindow strings.Builder
	)

	// HEADER / TABS
	header.WriteString(m.genTabs())
	header.WriteString(m.genTabHelp())

	if m.Debug {
		// One debug line
		header.WriteString(fmt.Sprintf("%s Width: %d Height: %d ErrorMsg: %s\n", styles.TextRed.Render("DEBUG ON:"), m.Globals.winW, m.Globals.winH, m.Globals.ErrorMsg))
	}

	if m.Globals.ErrorHelp != "" {
		m.Log.Println("Got error")
		header.WriteString(styles.ErrorHelp.Render(fmt.Sprintf("ERROR: %s", m.Globals.ErrorHelp)))
	} else {
		m.Log.Println("Got NO error, insert newline")
		//header.WriteString("\n")
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
		m.Log.Printf("CALL JobFromTemplate.View()\n")
		MainWindow.WriteString(m.JobFromTemplateTab.View(m.Log))

	case tabCluster:
		m.Log.Printf("CALL ClusterTab.View()\n")
		MainWindow.WriteString(m.ClusterTab.View(m.Log))

	case tabAbout:
		MainWindow.WriteString(m.tabAbout())
		// TODO: default
	}

	return lipgloss.JoinVertical(lipgloss.Left, header.String(), styles.MainWindow.Render(MainWindow.String()), styles.HelpWindow.Render(m.Help.View(keybindings.DefaultKeyMap)))
}
