package model

import (
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/lipgloss"
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
	case len(m.SacctSingleJobHist.Jobs) == 0:
		return fmt.Sprintf("Waiting for job %s info...\n", m.JobDetailsTab.SelJobID)

	}

	//width := m.Globals.winW - 10
	job := m.SacctSingleJobHist.Jobs[0]

	m.Log.Printf("Job Details req %#v ,got: %#v\n", m.JobDetailsTab.SelJobID, job.JobId)
	//scr = fmt.Sprintf("Job count: %d\n\n", len(m.SacctJob.Jobs))

	// TODO: consider moving this to a table...

	head := ""
	waitT := time.Unix(int64(*job.Time.Start), 0).Sub(time.Unix(int64(*job.Time.Submission), 0))
	runT := time.Unix(int64(*job.Time.End), 0).Sub(time.Unix(int64(*job.Time.Start), 0))
	fmtStr := "%-20s : %-40s\n"
	head += fmt.Sprintf(fmtStr, "Job ID", strconv.Itoa(*job.JobId))
	head += fmt.Sprintf(fmtStr, "User", *job.User)
	head += fmt.Sprintf(fmtStr, "Job Name", *job.Name)
	head += fmt.Sprintf(fmtStr, "Job Account", *job.Account)
	head += fmt.Sprintf(fmtStr, "Job Submission", time.Unix(int64(*job.Time.Submission), 0).String())
	head += fmt.Sprintf(fmtStr, "Job Start", time.Unix(int64(*job.Time.Start), 0).String())
	head += fmt.Sprintf(fmtStr, "Job End", time.Unix(int64(*job.Time.End), 0).String())
	head += fmt.Sprintf(fmtStr, "Job Wait time", waitT.String())
	head += fmt.Sprintf(fmtStr, "Job Run time", runT.String())
	head += fmt.Sprintf(fmtStr, "Partition", *job.Partition)
	head += fmt.Sprintf(fmtStr, "Priority", strconv.Itoa(*job.Priority))
	head += fmt.Sprintf(fmtStr, "QoS", *job.Qos)

	//scr += styles.JobStepBoxStyle.Width(width).Render(head)
	scr += styles.JobStepBoxStyle.Render(head)
	scr += fmt.Sprintf("\n Steps count: %d", len(*job.Steps))

	steps := ""
	for i, v := range *job.Steps {

		m.Log.Printf("Job Details, step: %d name: %s\n", i, *v.Step.Name)
		step := fmt.Sprintf(fmtStr, "Name", *v.Step.Name)
		step += fmt.Sprintf(fmtStr, "State", *v.State)
		step += fmt.Sprintf(fmtStr, "ExitStatus", *v.ExitCode.Status)
		if *v.ExitCode.Status == "SIGNALED" {
			step += fmt.Sprintf(fmtStr, "Signal ID", strconv.Itoa(*v.ExitCode.Signal.SignalId))
			step += fmt.Sprintf(fmtStr, "SignalName", *v.ExitCode.Signal.Name)
		}
		if v.KillRequestUser != nil {
			step += fmt.Sprintf(fmtStr, "KillReqUser", *v.KillRequestUser)
		}
		step += fmt.Sprintf(fmtStr, "Tasks", strconv.Itoa(*v.Tasks.Count))
		//steps += lipgloss.JoinVertical(lipgloss.Bottom, steps, styles.JobStepBoxStyle.Width(m.Globals.winW-10).Render(step))
		//
		//steps += "\n" + styles.JobStepBoxStyle.Width(width).Render(step)
		steps += "\n" + styles.JobStepBoxStyle.Render(step)
		//
		//m.Log.Printf("Step %d, VALUE: %q", i, steps)
		//m.Log.Printf("================================================================================\n")
	}
	scr += steps

	//// get Type of struct
	//t := reflect.TypeOf(m.JobDetailsTab.Jobs[0])
	//// get Value of struct
	//v := reflect.ValueOf(m.JobDetailsTab.Jobs[0])
	//// get struct fields from Type
	//f := reflect.VisibleFields(t)
	//scr += fmt.Sprintf("num. Fields = %d\n", len(f))
	//for _, val := range f {
	//	dv := v.FieldByName(val.Name).Elem()
	//	switch dv.Kind() {
	//	case reflect.String:
	//		scr += fmt.Sprintf("*string FIELD: %-20s VALUE: %-20s\n", val.Name, dv.String())
	//	case reflect.Int:
	//		scr += fmt.Sprintf("*int FIELD: %-20s VALUE: %-20d\n", val.Name, dv.Int())
	//	default:
	//		scr += fmt.Sprintf("UnknownKind %v FIELD: %-20s VALUE: %-20v\n", v.FieldByName(val.Name).Kind(), val.Name, v.FieldByName(val.Name).String())

	//	}
	//}
	//scr += fmt.Sprintf("Job:\n\n%#v\n\nSelected job: %#v\n\n", m.JobDetailsTab.SacctJob, m.JobDetailsTab.SelJobID)
	//m.LogF.WriteString(fmt.Sprintf("Job:\n\n%#v\n\nSelected job: %#v\n\n", m.JobDetailsTab.SacctJob, m.JobDetailsTab.SelJobID))

	return scr
	//return m.JobDetailsTab.SelJobID
}

func (m Model) tabJobFromTemplate() string {

	if m.EditTemplate {
		return m.TemplateEditor.View()
	} else {
		// TODO: if len(table)==0 return "no templates found"
		if len(m.JobFromTemplateTab.TemplatesList) == 0 {
			return styles.NotFound.Render("\nNo templates found!\n")
		} else {
			return m.TemplatesTable.View()
		}
	}
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

	s := "Version: " + version.BuildVersion + "\n"
	s += "Commit : " + version.BuildCommit + "\n"

	s += `
petar.jager@imba.oeaw.ac.at
CLIP-HPC Team @ VBC

`

	//st := styles.MainWindow.Copy().Height(m.Globals.winH - 10)

	//return "About tab active" + s
	//return styles.MainWindow.Render(styles.JobInfoBox.Render(s))
	return s
}

func (m Model) getJobInfo() string {
	var scr strings.Builder

	n := m.JobTab.SqueueTable.Cursor()
	m.Log.Printf("getJobInfo: cursor at %d table rows: %d\n", n, len(m.JobTab.SqueueFiltered.Jobs))
	if len(m.JobTab.SqueueFiltered.Jobs) == 0 || n == -1 {
		return "Select a job"
	}

	fmtStr := "%-15s : %-30s\n"
	fmtStrLast := "%-15s : %-30s"
	//ibFmt := "Job Name: %s\nJob Command: %s\nOutput: %s\nError: %s\n"
	infoBoxLeft := fmt.Sprintf(fmtStr, "Partition", *m.JobTab.SqueueFiltered.Jobs[n].Partition)
	infoBoxLeft += fmt.Sprintf(fmtStr, "QoS", *m.JobTab.SqueueFiltered.Jobs[n].Qos)
	infoBoxLeft += fmt.Sprintf(fmtStr, "TRES", *m.JobTab.SqueueFiltered.Jobs[n].TresReqStr)
	infoBoxLeft += fmt.Sprintf(fmtStr, "Batch Host", *m.JobTab.SqueueFiltered.Jobs[n].BatchHost)
	if m.JobTab.SqueueFiltered.Jobs[n].JobResources.Nodes != nil {
		infoBoxLeft += fmt.Sprintf(fmtStrLast, "AllocNodes", *m.JobTab.SqueueFiltered.Jobs[n].JobResources.Nodes)
	} else {
		infoBoxLeft += fmt.Sprintf(fmtStrLast, "AllocNodes", "none")

	}

	infoBoxRight := fmt.Sprintf(fmtStr, "Array Job ID", strconv.Itoa(*m.JobTab.SqueueFiltered.Jobs[n].ArrayJobId))
	if m.JobTab.SqueueFiltered.Jobs[n].ArrayTaskId != nil {
		infoBoxRight += fmt.Sprintf(fmtStr, "Array Task ID", strconv.Itoa(*m.JobTab.SqueueFiltered.Jobs[n].ArrayTaskId))
	} else {
		infoBoxRight += fmt.Sprintf(fmtStr, "Array Task ID", "NoTaskID")
	}
	infoBoxRight += fmt.Sprintf(fmtStr, "Gres Details", strings.Join(*m.JobTab.SqueueFiltered.Jobs[n].GresDetail, ","))
	infoBoxRight += fmt.Sprintf(fmtStr, "Features", *m.JobTab.SqueueFiltered.Jobs[n].Features)
	infoBoxRight += fmt.Sprintf(fmtStrLast, "wckey", *m.JobTab.SqueueFiltered.Jobs[n].Wckey)

	infoBoxMiddle := fmt.Sprintf(fmtStr, "Submit", time.Unix(*m.JobTab.SqueueFiltered.Jobs[n].SubmitTime, 0))
	if *m.JobTab.SqueueFiltered.Jobs[n].StartTime != 0 {
		infoBoxMiddle += fmt.Sprintf(fmtStrLast, "Start /expected", time.Unix(*m.JobTab.SqueueFiltered.Jobs[n].StartTime, 0))
	} else {
		infoBoxMiddle += fmt.Sprintf(fmtStrLast, "Start", "unknown")
	}

	infoBoxWide := fmt.Sprintf(fmtStr, "Job Name", *m.JobTab.SqueueFiltered.Jobs[n].Name)
	infoBoxWide += fmt.Sprintf(fmtStr, "Command", *m.JobTab.SqueueFiltered.Jobs[n].Command)
	infoBoxWide += fmt.Sprintf(fmtStr, "StdOut", *m.JobTab.SqueueFiltered.Jobs[n].StandardOutput)
	infoBoxWide += fmt.Sprintf(fmtStr, "StdErr", *m.JobTab.SqueueFiltered.Jobs[n].StandardError)
	infoBoxWide += fmt.Sprintf(fmtStrLast, "Working Dir", *m.JobTab.SqueueFiltered.Jobs[n].CurrentWorkingDirectory)

	// 8 for borders (~10 extra)
	//w := ((m.Globals.winW - 10) / 3) * 3
	//s := styles.JobInfoInBox.Copy().Width(w / 3).Height(5)
	////top := lipgloss.JoinHorizontal(lipgloss.Top, styles.JobInfoInBox.Render(infoBoxLeft), styles.JobInfoInBox.Render(infoBoxMiddle), styles.JobInfoInBox.Render(infoBoxRight))
	top := lipgloss.JoinHorizontal(lipgloss.Top, styles.JobInfoInBox.Render(infoBoxLeft), styles.JobInfoInBox.Render(infoBoxMiddle), styles.JobInfoInBox.Render(infoBoxRight))
	//s = styles.JobInfoInBox.Copy().Width(w + 4)
	scr.WriteString(lipgloss.JoinVertical(lipgloss.Left, top, styles.JobInfoInBottomBox.Render(infoBoxWide)))

	//return infoBox
	return scr.String()
}

func genTabHelp(t int) string {
	var th string
	switch t {
	case tabJobs:
		th = "List of jobs in the queue"
	case tabJobHist:
		th = "Last 7 days job history"
	case tabJobDetails:
		th = "Job details, select a job from Job History tab"
	case tabJobFromTemplate:
		th = "Edit and submit one of the job templates"
	case tabCluster:
		th = "List and status of cluster nodes"
	default:
		th = "SlurmCommander"
	}
	return th + "\n\n"
}

// Generate statistics string, horizontal.
func GenCountStr(cnt map[string]uint, l *log.Logger) string {
	var (
		scr string
	)

	sm := make([]struct {
		name string
		val  uint
	}, 0)

	// place map to slice
	for k, v := range cnt {
		sm = append(sm, struct {
			name string
			val  uint
		}{name: k, val: uint(v)})
	}

	// sort it
	sort.Slice(sm, func(i, j int) bool {
		if sm[i].name < sm[j].name {
			return true
		} else {
			return false
		}
	})

	// print it out
	scr = "Count: "
	for _, v := range sm {
		scr += fmt.Sprintf("%s: %d ", v.name, v.val)
	}
	scr += "\n\n"

	return scr
}

// Generate statistics string, vertical.
func GenCountStrVert(cnt map[string]uint, l *log.Logger) string {
	var (
		scr string
	)

	sm := make([]struct {
		name string
		val  uint
	}, 0)

	// place map to slice
	for k, v := range cnt {
		sm = append(sm, struct {
			name string
			val  uint
		}{name: k, val: uint(v)})
	}

	// sort it
	sort.Slice(sm, func(i, j int) bool {
		if sm[i].val > sm[j].val {
			return true
		} else {
			return false
		}
	})

	// print it out
	//scr = "Count: "
	for _, v := range sm {
		scr += fmt.Sprintf("%-15s: %d\n", v.name, v.val)
	}
	scr += "\n\n"

	return scr
}

func (m Model) JobClusterTabStats() string {

	m.Log.Printf("JobClusterTabStats called\n")

	str := "Nodes state statistics (filtered):\n\n"

	//str += GenCountStrVert(m.JobClusterTab.Stats.StateCnt, m.Log)
	str += GenCountStrVert(m.JobClusterTab.Stats.StateSimpleCnt, m.Log)

	return str
}

func HumanizeDuration(t time.Duration, l *log.Logger) string {
	var ret string

	l.Printf("Humanizing %f seconds to:\n", t.Seconds())
	// total seconds
	s := int64(t.Seconds())

	// days
	l.Printf("seconds: %d\n", s)
	d := s / (24 * 60 * 60)
	s = s % (24 * 60 * 60)
	l.Printf("seconds: %d\n", s)

	// hours
	h := s / 3600
	s = s % 3600
	l.Printf("seconds: %d\n", s)

	// minutes
	m := s / 60
	s = s % 60
	l.Printf("seconds: %d\n", s)

	ret += fmt.Sprintf("%.2d-%.2d:%.2d:%.2d", d, h, m, s)

	l.Printf("Humanized %q\n", ret)
	return ret
}

func (m Model) JobHistTabStats() string {

	m.Log.Printf("JobHistTabStats called\n")

	str := "History statistics (filtered):\n\n"
	str += GenCountStrVert(m.JobHistTab.Stats.StateCnt, m.Log)

	str += "Waiting times:\n\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinWait", HumanizeDuration(m.JobHistTab.Stats.MinWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "AvgWait", HumanizeDuration(m.JobHistTab.Stats.AvgWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MedWait", HumanizeDuration(m.JobHistTab.Stats.MedWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MaxWait", HumanizeDuration(m.JobHistTab.Stats.MaxWait, m.Log))

	return str
}

func (m Model) JobTabStats() string {

	m.Log.Printf("JobTabStats called\n")

	str := "Queue statistics (filtered):\n\n"
	// TODO: make it sorted
	str += GenCountStrVert(m.JobTab.Stats.StateCnt, m.Log)
	//for k, v := range m.JobTab.Stats.StateCnt {
	//	str += fmt.Sprintf("%-10s : %d\n", k, v)
	//}
	str += "Pending Jobs:\n\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinWait", HumanizeDuration(m.JobTab.Stats.MinWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "AvgWait", HumanizeDuration(m.JobTab.Stats.AvgWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MedWait", HumanizeDuration(m.JobTab.Stats.MedWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MaxWait", HumanizeDuration(m.JobTab.Stats.MaxWait, m.Log))

	str += "\nRunning Jobs:\n\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinRun", HumanizeDuration(m.JobTab.Stats.MinRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "AvgRun", HumanizeDuration(m.JobTab.Stats.AvgRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MedRun", HumanizeDuration(m.JobTab.Stats.MedRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MaxRun", HumanizeDuration(m.JobTab.Stats.MaxRun, m.Log))

	return str
}

func (m Model) View() string {

	var (
		scr        strings.Builder
		MainWindow strings.Builder
	)

	// HEADER / TABS
	scr.WriteString(m.genTabs())
	scr.WriteString(genTabHelp(int(m.ActiveTab)))

	// PICK and RENDER ACTIVE TAB

	switch m.ActiveTab {
	case tabJobs:
		// Top Main
		MainWindow.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n", m.JobTab.Filter.Value(), len(m.JobTab.SqueueFiltered.Jobs)))
		MainWindow.WriteString(GenCountStr(m.JobTab.Stats.StateCnt, m.Log))

		// Mid Main: table || table+stats || table+menu
		switch {
		case m.JobTab.MenuOn:
			// table + menu
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobs(), styles.MenuBoxStyle.Render(m.JobTab.Menu.View())))
			m.Log.Printf("\nITEMS LIST: %#v\n", m.JobTab.Menu.Items())
		case m.JobTab.StatsOn:
			// table + stats
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobs(), styles.MenuBoxStyle.Render(m.JobTabStats())))
		default:
			// table
			MainWindow.WriteString(m.tabJobs())
		}

		// Low Main: nil || info || filter
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			// filter
			MainWindow.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		case m.JobTab.InfoOn:
			// info
			MainWindow.WriteString(styles.JobInfoBox.Render(m.getJobInfo()))
		}
	case tabJobHist:
		// Top Main
		MainWindow.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n", m.JobHistTab.Filter.Value(), len(m.JobHistTab.SacctHistFiltered.Jobs)))
		MainWindow.WriteString(GenCountStr(m.JobHistTab.Stats.StateCnt, m.Log))
		MainWindow.WriteString(fmt.Sprintf("AvgWait: %s MedianWait: %s MinWait: %s Maxwait: %s\n", m.JobHistTab.Stats.AvgWait.String(), m.JobHistTab.MedWait.String(), m.JobHistTab.MinWait.String(), m.JobHistTab.MaxWait.String()))

		// Mid Main: table || table+stats
		switch {
		case m.JobHistTab.StatsOn:
			// table + stats
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobHist(), styles.MenuBoxStyle.Render(m.JobHistTabStats())))
		default:
			// table
			MainWindow.WriteString(m.tabJobHist())
		}

		// Low Main: nil || filter
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			// filter
			MainWindow.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobHistTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		}
	case tabJobDetails:
		MainWindow.WriteString(m.tabJobDetails())
	case tabJobFromTemplate:
		MainWindow.WriteString(m.tabJobFromTemplate())
	case tabCluster:
		// Top Main
		MainWindow.WriteString(fmt.Sprintf("Filter: %10.10s\tItems: %d\n\n", m.JobClusterTab.Filter.Value(), len(m.JobClusterTab.SinfoFiltered.Nodes)))
		MainWindow.WriteString(GenCountStr(m.JobClusterTab.Stats.StateCnt, m.Log))

		// Mid Main: table || table+stats
		switch {
		case m.JobClusterTab.StatsOn:
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabCluster(), styles.MenuBoxStyle.Render(m.JobClusterTabStats())))
		default:
			MainWindow.WriteString(m.tabCluster())
		}

		// Low Main: nil || filter
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			// filter
			MainWindow.WriteString(fmt.Sprintf("Filter value (search accross all fields!):\n%s\n%s", m.JobClusterTab.Filter.View(), "(Enter to finish, Esc to clear filter and abort)") + "\n")
		}
	case tabAbout:
		//MainWindow.WriteString(m.tabAbout())
		MainWindow.WriteString(m.tabAbout())
	}

	// FOOTER
	//MainWindow.WriteString("\n")
	// Debug information:
	//if m.Globals.Debug {
	//	scr.WriteString("DEBUG:\n")
	//	scr.WriteString(fmt.Sprintf("Last key pressed: %q\n", m.lastKey))
	//	scr.WriteString(fmt.Sprintf("Window Width: %d\tHeight:%d\n", m.winW, m.winH))
	//	scr.WriteString(fmt.Sprintf("Active tab: %d\t Active Filter value: TBD\t InfoOn: %v\n", m.ActiveTab, m.InfoOn))
	//	scr.WriteString(fmt.Sprintf("Debug Msg: %q\n", m.DebugMsg))
	//}

	// TODO: Help doesn't split into multiple lines (e.g. when window too narrow)
	//scr.WriteString(m.Help.View(keybindings.DefaultKeyMap))

	//scr.WriteString(styles.MainWindow.Render(MainWindow.String()))
	//scr.WriteString(styles.HelpWindow.Render(m.Help.View(keybindings.DefaultKeyMap)))
	scr.WriteString(lipgloss.JoinVertical(lipgloss.Left, styles.MainWindow.Render(MainWindow.String()), styles.HelpWindow.Render(m.Help.View(keybindings.DefaultKeyMap))))

	return scr.String()
}
