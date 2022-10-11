package model

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
	"github.com/pja237/slurmcommander/internal/keybindings"
	"github.com/pja237/slurmcommander/internal/styles"
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

	//row := lipgloss.JoinHorizontal(
	//	lipgloss.Top,
	//	activeTab.Render("Jobs"),
	//	tab.Render("Job History"),
	//	tab.Render("Cluster"),
	//	tab.Render("About"),
	//)

	//gap := tabGap.Render(strings.Repeat(" ", max(0, width-lipgloss.Width(row)-2)))
	gap := styles.TabGap.Render(strings.Repeat(" ", max(0, m.winW-lipgloss.Width(row)-2)))
	row = lipgloss.JoinHorizontal(lipgloss.Bottom, row, gap)
	//doc.WriteString(row + "\n\n")
	doc.WriteString(row + "\n")

	return doc.String()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (m Model) tabJobs() string {

	// TODO: See what's more visually clear way to present info
	// e.g. Show selected job info in:
	// a) separate toggle-able window on the side OR
	// b) header/footer (above/below) of the table, like `top` does
	return m.SqueueTable.View() + "\n"
}

func (m Model) tabJobHist() string {

	// TODO: do some statistics on job history
	// e.g. avg waiting times, jobs successfull/failed count, etc...

	return m.JobHistTab.SacctTable.View() + "\n"
}

func (m Model) tabJobDetails() (scr string) {

	// race between View() call and command.SingleJobGetSacct(m.JobDetailsTab.SelJobID) call
	switch {
	case m.JobDetailsTab.SelJobID == "":
		return "Select a job from the Job History tab.\n"
	case len(m.SacctJob.Jobs) == 0:
		return fmt.Sprintf("Waiting for job %s info...\n", m.JobDetailsTab.SelJobID)

	}

	m.Log.Printf("Job Account: %#v\n", *m.SacctJob.Jobs[0].Account)
	scr = fmt.Sprintf("Job count: %d\n\n", len(m.SacctJob.Jobs))
	//for i, v := range m.SacctJob.Jobs {
	//	scr += fmt.Sprintf("Job: %d\n\n%#v\n\nSelected job: %#v\n\n", i, v, m.JobDetailsTab.SelJobID)
	//}

	// TODO: consider moving this to a table...

	waitT := time.Unix(int64(*m.SacctJob.Jobs[0].Time.Submission), 0).Sub(time.Unix(int64(*m.SacctJob.Jobs[0].Time.Submission), 0))
	runT := time.Unix(int64(*m.SacctJob.Jobs[0].Time.End), 0).Sub(time.Unix(int64(*m.SacctJob.Jobs[0].Time.Start), 0))
	fmtStr := "%-20s : %-40s\n"
	scr += fmt.Sprintf(fmtStr, "Job ID", strconv.Itoa(*m.SacctJob.Jobs[0].JobId))
	scr += fmt.Sprintf(fmtStr, "Job Name", *m.SacctJob.Jobs[0].Name)
	scr += fmt.Sprintf(fmtStr, "Job Account", *m.SacctJob.Jobs[0].Account)
	scr += fmt.Sprintf(fmtStr, "Job Submission", time.Unix(int64(*m.SacctJob.Jobs[0].Time.Submission), 0).String())
	scr += fmt.Sprintf(fmtStr, "Job Start", time.Unix(int64(*m.SacctJob.Jobs[0].Time.Start), 0).String())
	scr += fmt.Sprintf(fmtStr, "Job End", time.Unix(int64(*m.SacctJob.Jobs[0].Time.End), 0).String())
	scr += fmt.Sprintf(fmtStr, "Job Wait time", waitT.String())
	scr += fmt.Sprintf(fmtStr, "Job Run time", runT.String())
	scr += fmt.Sprintf(fmtStr, "Partition", *m.SacctJob.Jobs[0].Partition)
	scr += fmt.Sprintf(fmtStr, "Priority", strconv.Itoa(*m.SacctJob.Jobs[0].Priority))
	scr += fmt.Sprintf(fmtStr, "QoS", *m.SacctJob.Jobs[0].Qos)
	scr += "---\n"
	scr += fmt.Sprintf("Job:\n\n%#v\n\nSelected job: %#v\n\n", m.JobDetailsTab.SacctJob, m.JobDetailsTab.SelJobID)
	//m.LogF.WriteString(fmt.Sprintf("Job:\n\n%#v\n\nSelected job: %#v\n\n", m.JobDetailsTab.SacctJob, m.JobDetailsTab.SelJobID))

	return scr
	//return m.JobDetailsTab.SelJobID
}

func (m Model) tabJobFromTemplate() string {

	if m.EditTemplate {
		//return fmt.Sprintf("%s\n\n", m.TemplateEditor.Placeholder) + m.TemplateEditor.View()
		return m.TemplateEditor.View()
	} else {
		return m.TemplatesTable.View()
	}
	//return "Jobs from Template tab active"
}

func (m Model) tabCluster() string {

	var (
		scr     string = ""
		cpuPerc float64
		memPerc float64
	)

	// node info
	// TODO: rework, doesn't work when table filtering is on
	sel := m.SinfoTable.Cursor()
	m.Log.Printf("ClusterTab Selected: %d\n", sel)
	m.Log.Printf("ClusterTab len results: %d\n", len(m.JobClusterTab.SinfoFiltered.Nodes))
	if len(m.JobClusterTab.SinfoFiltered.Nodes) > 0 && sel != -1 {
		m.CpuBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
		cpuPerc = float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].AllocCpus) / float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].Cpus)
		//m.CpuBar.SetPercent(cpuPerc)
		m.MemBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
		memPerc = float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].AllocMemory) / float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].RealMemory)
	} else {
		cpuPerc = 0
		memPerc = 0
	}

	scr += "Cpu and memory utilization:\n"
	scr += fmt.Sprintf("cpuPerc: %.2f ", cpuPerc)
	scr += m.CpuBar.ViewAs(cpuPerc)
	scr += "\n"
	scr += fmt.Sprintf("memPerc: %.2f ", memPerc)
	scr += m.MemBar.ViewAs(memPerc)
	scr += "\n\n"

	// table
	scr += m.SinfoTable.View() + "\n"

	return scr
	//return m.SinfoTable.View()

	//return "Cluster tab active"
	// TODO: HEADER/FOOTER that shows details for selected node
	// e.g. progress bars with cpu/mem usage (percentages)
	// TODO: reselect table columns (move mem/cpu to header/footer above and pick others, e.g. partition? features? think... )
	// TODO: rename to "Cluster Nodes", add "Cluster QoS/Partition" tab? OR find an elegant way to group those in one tab?
}

func (m Model) tabAbout() string {

	s := `

petar.jager@imba.oeaw.ac.at
CLIP-HPC Team @ VBC
	`

	return "About tab active" + s
}

func (m Model) getJobInfo() string {
	// TODO:
	// fix: if after filtering m.table.Cursor|SelectedRow > lines in table, Info crashes trying to fetch nonexistent row
	//return strconv.Itoa(m.SqueueTable.Cursor()) + "\n" + m.JobTab.SelectedJob + "\n" + m.JobTab.MenuChoice.Title()
	n := m.JobTab.SqueueTable.Cursor()
	ibFmt := "Job Name: %s\nJob Command: %s\nOutput: %s\nError: %s\n"
	infoBox := fmt.Sprintf(ibFmt, *m.JobTab.SqueueFiltered.Jobs[n].Name,
		*m.JobTab.SqueueFiltered.Jobs[n].Command,
		*m.JobTab.SqueueFiltered.Jobs[n].StandardOutput,
		*m.JobTab.SqueueFiltered.Jobs[n].StandardError)
	//infoBox := strconv.Itoa(m.SqueueTable.Cursor()) + "\n" +
	//	m.JobTab.SelectedJob + "\n" +
	//	m.JobTab.MenuChoice.Title() + "\n" +
	//	m.JobTab.SqueueTable.SelectedRow()[0] + "\n" +
	//	*m.JobTab.SqueueFiltered.Jobs[n].Name
	return infoBox
}

func genTabHelp(t int) string {
	var th string
	switch t {
	case tabJobs:
		th = "Job queue list"
	default:
		th = "default tab help"
	}
	return th + "\n"
}

func (m Model) View() string {

	var scr strings.Builder

	// HEADER / TABS
	scr.WriteString(m.genTabs())
	scr.WriteString(genTabHelp(int(m.ActiveTab)))

	// PICK and RENDER ACTIVE TAB

	switch m.ActiveTab {
	case tabJobs:
		//scr.WriteString("Filter: " + m.JobTab.Filter.Value() + "\n\n")
		scr.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n\n", m.JobTab.Filter.Value(), len(m.JobTab.SqueueFiltered.Jobs)))
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			scr.WriteString(m.tabJobs())
			scr.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		case m.JobTab.MenuOn:
			// TODO: Render menu here
			scr.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobs(), styles.JobInfoBox.Render(m.JobTab.Menu.View())))
			//scr.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobs(), m.JobTab.Menu.View()))
			m.Log.Printf("\nITEMS LIST: %#v\n", m.JobTab.Menu.Items())
		case m.JobTab.InfoOn:
			//scr.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobs(), focusedModelStyle.Render(m.getJobInfo())))
			scr.WriteString(m.tabJobs() + "\n")
			scr.WriteString(styles.JobInfoBox.Render(m.getJobInfo()))
		default:
			scr.WriteString(m.tabJobs())
		}
	case tabJobHist:
		//scr.WriteString("Filter: " + m.JobHistTab.Filter.Value() + "\n\n")
		scr.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n\n", m.JobHistTab.Filter.Value(), len(m.JobHistTab.SacctListFiltered)))
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			scr.WriteString(m.tabJobHist())
			scr.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobHistTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		default:
			scr.WriteString(m.tabJobHist())
		}
	case tabJobDetails:
		scr.WriteString(m.tabJobDetails())
	case tabJobFromTemplate:
		scr.WriteString(m.tabJobFromTemplate())
	case tabCluster:
		//scr.WriteString("Filter: " + m.JobClusterTab.Filter.Value() + "\n\n")
		scr.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n\n", m.JobClusterTab.Filter.Value(), len(m.JobClusterTab.SinfoFiltered.Nodes)))
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			scr.WriteString(m.tabCluster())
			scr.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobClusterTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		default:
			scr.WriteString(m.tabCluster())
		}
	case tabAbout:
		scr.WriteString(m.tabAbout())
	}

	// FOOTER
	scr.WriteString("\n")
	// Debug information:
	if m.Globals.Debug {
		scr.WriteString("DEBUG:\n")
		scr.WriteString(fmt.Sprintf("Last key pressed: %q\n", m.lastKey))
		scr.WriteString(fmt.Sprintf("Window Width: %d\tHeight:%d\n", m.winW, m.winH))
		scr.WriteString(fmt.Sprintf("Active tab: %d\t Active Filter value: TBD\t InfoOn: %v\n", m.ActiveTab, m.InfoOn))
		scr.WriteString(fmt.Sprintf("Debug Msg: %q\n", m.DebugMsg))
	}

	// TODO: Help doesn't split into multiple lines (e.g. when window too narrow)
	scr.WriteString(m.Help.View(keybindings.DefaultKeyMap))

	return scr.String()
}
