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

	return m.SqueueTable.View() + "\n"
}

func (m Model) tabJobHist() string {

	return m.JobHistTab.SacctTable.View() + "\n"
}

func (m Model) tabJobDetails() (scr string) {

	var (
		runT  time.Duration
		waitT time.Duration
	)

	// race between View() call and command.SingleJobGetSacct(m.JobDetailsTab.SelJobID) call
	switch {
	//case m.JobDetailsTab.SelJobID == "":
	//	return "Select a job from the Job History tab.\n"
	case m.JobDetailsTab.SelJobIDNew == -1:
		return "Select a job from the Job History tab.\n"
	//case len(m.SacctSingleJobHist.Jobs) == 0:
	//	return fmt.Sprintf("Waiting for job %s info...\n", m.JobDetailsTab.SelJobID)
	case len(m.JobHistTab.SacctHistFiltered.Jobs) == 0:
		//return fmt.Sprintf("Waiting for job %s info...\n", m.JobDetailsTab.SelJobID)
		return "Select a job from the Job History tab.\n"
	}

	//width := m.Globals.winW - 10

	//job := m.SacctSingleJobHist.Jobs[0]
	// NEW:
	job := m.JobHistTab.SacctHistFiltered.Jobs[m.JobDetailsTab.SelJobIDNew]

	m.Log.Printf("Job Details req %#v ,got: %#v\n", m.JobDetailsTab.SelJobID, job.JobId)

	// TODO: consider moving this to a viewport...

	fmtStr := "%-20s : %-60s\n"
	fmtStrX := "%-20s : %-60s"

	head := ""
	waitT = time.Unix(int64(*job.Time.Start), 0).Sub(time.Unix(int64(*job.Time.Submission), 0))
	// If job is RUNNING, use Elapsed instead of Sub (because End=0)
	if *job.State.Current == "RUNNING" {
		runT = time.Duration(int64(*job.Time.Elapsed) * int64(time.Second))
	} else {
		runT = time.Unix(int64(*job.Time.End), 0).Sub(time.Unix(int64(*job.Time.Start), 0))
	}

	head += styles.StatsSeparatorTitle.Render(fmt.Sprintf(fmtStrX, "Job ID", strconv.Itoa(*job.JobId)))
	head += "\n"
	head += fmt.Sprintf(fmtStr, "Job Name", *job.Name)
	head += fmt.Sprintf(fmtStr, "User", *job.User)
	head += fmt.Sprintf(fmtStr, "Group", *job.Group)
	head += fmt.Sprintf(fmtStr, "Job Account", *job.Account)
	head += fmt.Sprintf(fmtStr, "Job Submission", time.Unix(int64(*job.Time.Submission), 0).String())
	head += fmt.Sprintf(fmtStr, "Job Start", time.Unix(int64(*job.Time.Start), 0).String())
	// Running jobs have End==0
	if *job.State.Current == "RUNNING" {
		head += fmt.Sprintf(fmtStr, "Job End", "RUNNING")
	} else {
		head += fmt.Sprintf(fmtStr, "Job End", time.Unix(int64(*job.Time.End), 0).String())
	}
	head += fmt.Sprintf(fmtStr, "Job Wait time", waitT.String())
	head += fmt.Sprintf(fmtStr, "Job Run time", runT.String())
	head += fmt.Sprintf(fmtStr, "Partition", *job.Partition)
	head += fmt.Sprintf(fmtStr, "Priority", strconv.Itoa(*job.Priority))
	head += fmt.Sprintf(fmtStr, "QoS", *job.Qos)

	scr += styles.JobStepBoxStyle.Width(90).Render(head)
	scr += "\n"

	scr += styles.TextYellow.Render(fmt.Sprintf("Steps count: %d", len(*job.Steps)))

	steps := ""
	for i, v := range *job.Steps {

		m.Log.Printf("Job Details, step: %d name: %s\n", i, *v.Step.Name)
		step := styles.StatsSeparatorTitle.Render(fmt.Sprintf(fmtStrX, "Name", *v.Step.Name))
		step += "\n"
		step += fmt.Sprintf(fmtStr, "Nodes", *v.Nodes.Range)
		if *v.State != "COMPLETED" {
			step += styles.JobStepExitStatusRed.Render(fmt.Sprintf(fmtStrX, "State", *v.State))
			step += "\n"
		} else {
			//step += fmt.Sprintf(fmtStr, "State", *v.State)
			step += styles.JobStepExitStatusGreen.Render(fmt.Sprintf(fmtStrX, "State", *v.State))
			step += "\n"
		}
		if *v.ExitCode.Status != "SUCCESS" {
			step += styles.JobStepExitStatusRed.Render(fmt.Sprintf(fmtStrX, "ExitStatus", *v.ExitCode.Status))
			step += "\n"
		} else {
			step += styles.JobStepExitStatusGreen.Render(fmt.Sprintf(fmtStrX, "ExitStatus", *v.ExitCode.Status))
			step += "\n"
		}
		if *v.ExitCode.Status == "SIGNALED" {
			step += styles.JobStepExitStatusRed.Render(fmt.Sprintf(fmtStrX, "Signal ID", strconv.Itoa(*v.ExitCode.Signal.SignalId)))
			step += "\n"
			step += styles.JobStepExitStatusRed.Render(fmt.Sprintf(fmtStrX, "SignalName", *v.ExitCode.Signal.Name))
			step += "\n"
		}
		if v.KillRequestUser != nil {
			step += fmt.Sprintf(fmtStr, "KillReqUser", *v.KillRequestUser)
		}
		step += fmt.Sprintf(fmtStr, "Tasks", strconv.Itoa(*v.Tasks.Count))

		// TODO: TRES part needs quite some love...
		tres := ""
		tresAlloc := ""

		//tresReqMin := ""
		//tresReqMax := ""
		//tresReqAvg := ""
		//tresReqTotal := ""
		//tresConMax := ""
		//tresConMin := ""
		// TRES: allocated
		tresAlloc += "\nALLOCATED:\n"
		m.Log.Printf("Dumping step allocation: %#v\n", *v.Tres.Allocated)
		m.Log.Printf("ALLOCATED:\n")
		for i, t := range *v.Tres.Allocated {
			if t.Count != nil {
				m.Log.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
				tresAlloc += "* "
				if *t.Type == "gres" {
					// TODO:
					//fmtStr := "%-20s : %-60s\n"
					tresAlloc += fmt.Sprintf(fmtStr, *t.Type, strings.Join([]string{*t.Name, strconv.Itoa(*t.Count)}, ":"))
				} else {
					// TODO:
					tresAlloc += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
				}
			}
		}
		//// REQUESTED:MIN
		//tresReqMin += "REQUESTED:Min:\n"
		//m.Log.Printf("REQ:Min\n")
		//for i, t := range *v.Tres.Requested.Min {
		//	if t.Count != nil {
		//		m.Log.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresReqMin += " "
		//		tresReqMin += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// REQUESTED:MAX
		//m.Log.Printf("REQ:Max\n")
		//tresReqMax += "REQUESTED:Max:\n"
		//for i, t := range *v.Tres.Requested.Min {
		//	if t.Count != nil {
		//		m.Log.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresReqMax += " "
		//		tresReqMax += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// REQUESTED:AVG
		//m.Log.Printf("REQ:Avg\n")
		//tresReqAvg += "REQUESTED:Avg:\n"
		//for i, t := range *v.Tres.Requested.Average {
		//	if t.Count != nil {
		//		m.Log.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresReqAvg += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// REQUESTED:TOT
		//tresReqAvg += "REQUESTED:Tot:\n"
		//m.Log.Printf("REQ:Tot\n")
		//for i, t := range *v.Tres.Requested.Total {
		//	if t.Count != nil {
		//		m.Log.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresReqTotal += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// Consumed:Min
		//tresConMin += "CONSUMED:Min:\n"
		//m.Log.Printf("CONS:Min\n")
		//for i, t := range *v.Tres.Consumed.Min {
		//	if t.Count != nil {
		//		m.Log.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresConMin += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// Consumed:Max
		//tresConMax += "CONSUMED:Max:\n"
		//m.Log.Printf("CONS:Max\n")
		//for i, t := range *v.Tres.Consumed.Max {
		//	if t.Count != nil {
		//		m.Log.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresConMax += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//tres = lipgloss.JoinHorizontal(lipgloss.Top, styles.TresBox.Render(tresAlloc), styles.TresBox.Width(40).Render(tresConMax))

		// For now, show just allocated, later rework this whole part
		tres = styles.TresBox.Render(tresAlloc)

		step += tres

		// when the step is finished, append it to steps string
		steps += "\n" + styles.JobStepBoxStyle.Render(step)
	}
	scr += steps

	return scr
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

func (m Model) tabClusterBars() string {
	var (
		scr     string = ""
		cpuPerc float64
		memPerc float64
	)

	sel := m.SinfoTable.Cursor()
	m.Log.Printf("ClusterTab Selected: %d\n", sel)
	m.Log.Printf("ClusterTab len results: %d\n", len(m.JobClusterTab.SinfoFiltered.Nodes))
	m.JobClusterTab.CpuBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
	m.JobClusterTab.MemBar = progress.New(progress.WithGradient("#277BC0", "#FFCB42"))
	if len(m.JobClusterTab.SinfoFiltered.Nodes) > 0 && sel != -1 {
		cpuPerc = float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].AllocCpus) / float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].Cpus)
		memPerc = float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].AllocMemory) / float64(*m.JobClusterTab.SinfoFiltered.Nodes[sel].RealMemory)

		scr += fmt.Sprintf("CPU used/total: %d/%d\n", *m.JobClusterTab.SinfoFiltered.Nodes[sel].AllocCpus, *m.JobClusterTab.SinfoFiltered.Nodes[sel].Cpus)
		scr += m.CpuBar.ViewAs(cpuPerc)
		scr += "\n"
		scr += fmt.Sprintf("MEM used/total: %d/%d\n", *m.JobClusterTab.SinfoFiltered.Nodes[sel].AllocMemory, *m.JobClusterTab.SinfoFiltered.Nodes[sel].RealMemory)
		scr += m.MemBar.ViewAs(memPerc)
		scr += "\n\n"
	} else {
		cpuPerc = 0
		memPerc = 0
		scr += fmt.Sprintf("CPU used/total: %d/%d\n", 0, 0)
		scr += m.CpuBar.ViewAs(cpuPerc)
		scr += "\n"
		scr += fmt.Sprintf("MEM used/total: %d/%d\n", 0, 0)
		scr += m.MemBar.ViewAs(memPerc)
		scr += "\n\n"

	}

	return scr
}
func (m Model) tabCluster() string {

	scr := m.SinfoTable.View() + "\n"

	return scr
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

func (m Model) getJobHistCounts() string {
	var (
		ret   string
		top5u string
		top5a string
		jpp   string
		jpq   string
	)

	fmtStr := "%-20s : %6d\n"
	fmtTitle := "%-29s"

	top5u += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Top 5 User"))
	top5u += "\n"
	for _, v := range m.JobHistTab.Breakdowns.Top5user {
		top5u += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	top5a += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Top 5 Accounts"))
	top5a += "\n"
	for _, v := range m.JobHistTab.Breakdowns.Top5acc {
		top5a += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	jpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Jobs per Partition"))
	jpp += "\n"
	for _, v := range m.JobHistTab.Breakdowns.JobPerPart {
		jpp += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	jpq += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Jobs per QoS"))
	jpq += "\n"
	for _, v := range m.JobHistTab.Breakdowns.JobPerQos {
		jpq += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	top5u = styles.CountsBox.Render(top5u)
	top5a = styles.CountsBox.Render(top5a)
	jpq = styles.CountsBox.Render(jpq)
	jpp = styles.CountsBox.Render(jpp)

	ret = lipgloss.JoinHorizontal(lipgloss.Top, top5u, top5a, jpp, jpq)

	return ret
}

func (m Model) getJobCounts() string {
	var (
		ret   string
		top5u string
		top5a string
		jpp   string
		jpq   string
	)

	fmtStr := "%-20s : %6d\n"
	fmtTitle := "%-29s"

	top5u += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Top 5 User"))
	top5u += "\n"
	for _, v := range m.JobTab.Breakdowns.Top5user {
		top5u += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	top5a += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Top 5 Accounts"))
	top5a += "\n"
	for _, v := range m.JobTab.Breakdowns.Top5acc {
		top5a += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	jpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Jobs per Partition"))
	jpp += "\n"
	for _, v := range m.JobTab.Breakdowns.JobPerPart {
		jpp += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	jpq += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Jobs per QoS"))
	jpq += "\n"
	for _, v := range m.JobTab.Breakdowns.JobPerQos {
		jpq += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	top5u = styles.CountsBox.Render(top5u)
	top5a = styles.CountsBox.Render(top5a)
	jpq = styles.CountsBox.Render(jpq)
	jpp = styles.CountsBox.Render(jpp)

	ret = lipgloss.JoinHorizontal(lipgloss.Top, top5u, top5a, jpp, jpq)

	return ret
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
		infoBoxMiddle += fmt.Sprintf(fmtStr, "Start /expected", time.Unix(*m.JobTab.SqueueFiltered.Jobs[n].StartTime, 0))
	} else {
		infoBoxMiddle += fmt.Sprintf(fmtStr, "Start", "unknown")
	}
	// placeholder lines
	infoBoxMiddle += "\n"
	infoBoxMiddle += "\n"
	// EO placeholder lines
	infoBoxMiddle += fmt.Sprintf(fmtStrLast, "State reason", *m.JobTab.SqueueFiltered.Jobs[n].StateReason)

	infoBoxWide := fmt.Sprintf(fmtStr, "Job Name", *m.JobTab.SqueueFiltered.Jobs[n].Name)
	infoBoxWide += fmt.Sprintf(fmtStr, "Command", *m.JobTab.SqueueFiltered.Jobs[n].Command)
	infoBoxWide += fmt.Sprintf(fmtStr, "StdOut", *m.JobTab.SqueueFiltered.Jobs[n].StandardOutput)
	infoBoxWide += fmt.Sprintf(fmtStr, "StdErr", *m.JobTab.SqueueFiltered.Jobs[n].StandardError)
	infoBoxWide += fmt.Sprintf(fmtStrLast, "Working Dir", *m.JobTab.SqueueFiltered.Jobs[n].CurrentWorkingDirectory)

	// 8 for borders (~10 extra)
	//w := ((m.Globals.winW - 10) / 3) * 3
	//s := styles.JobInfoInBox.Copy().Width(w / 3).Height(5)
	////top := lipgloss.JoinHorizontal(lipgloss.Top, styles.JobInfoInBox.Render(infoBoxLeft), styles.JobInfoInBox.Render(infoBoxMiddle), styles.JobInfoInBox.Render(infoBoxRight))
	// TODO: use builder here
	top := lipgloss.JoinHorizontal(lipgloss.Top, styles.JobInfoInBox.Render(infoBoxLeft), styles.JobInfoInBox.Render(infoBoxMiddle), styles.JobInfoInBox.Render(infoBoxRight))
	//s = styles.JobInfoInBox.Copy().Width(w + 4)
	scr.WriteString(lipgloss.JoinVertical(lipgloss.Left, top, styles.JobInfoInBottomBox.Render(infoBoxWide)))

	//return infoBox
	return scr.String()
}

func (m *Model) genTabHelp() string {
	var th string
	switch m.ActiveTab {
	case tabJobs:
		th = "List of jobs in the queue"
	case tabJobHist:
		th = fmt.Sprintf("List of jobs in the last %d days from all user associated accounts. (timeout: %d seconds)", m.Globals.JobHistStart, m.Globals.JobHistTimeout)
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

	// sort first by name
	sort.Slice(sm, func(i, j int) bool {
		if sm[i].name < sm[j].name {
			return true
		} else {
			return false
		}
	})
	// then sort by numbers
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
	var str string

	m.Log.Printf("JobClusterTabStats called\n")

	sel := m.JobClusterTab.SinfoTable.Cursor()
	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Nodes states (filtered):"))
	str += "\n"

	if len(m.JobClusterTab.SinfoFiltered.Nodes) > 0 {
		//str += GenCountStrVert(m.JobClusterTab.Stats.StateCnt, m.Log)
		str += GenCountStrVert(m.JobClusterTab.Stats.StateSimpleCnt, m.Log)
	}

	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Selected node:"))

	if len(m.JobClusterTab.SinfoFiltered.Nodes) > 0 && sel != -1 {
		str += "\n"
		str += fmt.Sprintf("%-15s: %s\n", "Arch", *m.JobClusterTab.SinfoFiltered.Nodes[sel].Architecture)
		str += fmt.Sprintf("%-15s: %s\n", "Features", *m.JobClusterTab.SinfoFiltered.Nodes[sel].ActiveFeatures)
		str += fmt.Sprintf("%-15s: %s\n", "TRES", *m.JobClusterTab.SinfoFiltered.Nodes[sel].Tres)
		if m.JobClusterTab.SinfoFiltered.Nodes[sel].TresUsed != nil {
			str += fmt.Sprintf("%-15s: %s\n", "TRES Used", *m.JobClusterTab.SinfoFiltered.Nodes[sel].TresUsed)
		} else {
			str += fmt.Sprintf("%-15s: %s\n", "TRES Used", "")
		}
		str += fmt.Sprintf("%-15s: %s\n", "GRES", *m.JobClusterTab.SinfoFiltered.Nodes[sel].Gres)
		str += fmt.Sprintf("%-15s: %s\n", "GRES Used", *m.JobClusterTab.SinfoFiltered.Nodes[sel].GresUsed)
		str += fmt.Sprintf("%-15s: %s\n", "Partitions", strings.Join(*m.JobClusterTab.SinfoFiltered.Nodes[sel].Partitions, ","))
	}
	return str
}

func HumanizeDuration(t time.Duration, l *log.Logger) string {
	var ret string

	// total seconds
	s := int64(t.Seconds())

	// days
	d := s / (24 * 60 * 60)
	s = s % (24 * 60 * 60)

	// hours
	h := s / 3600
	s = s % 3600

	// minutes
	m := s / 60
	s = s % 60

	ret += fmt.Sprintf("%.2d-%.2d:%.2d:%.2d", d, h, m, s)

	l.Printf("Humanized %f to %q\n", t.Seconds(), ret)
	return ret
}

func (m Model) JobHistTabStats() string {

	m.Log.Printf("JobHistTabStats called\n")

	str := styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Historical job states (filtered):"))
	str += "\n\n"
	str += GenCountStrVert(m.JobHistTab.Stats.StateCnt, m.Log)

	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Waiting times (finished jobs):"))
	str += "\n\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinWait", HumanizeDuration(m.JobHistTab.Stats.MinWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "AvgWait", HumanizeDuration(m.JobHistTab.Stats.AvgWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MedWait", HumanizeDuration(m.JobHistTab.Stats.MedWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MaxWait", HumanizeDuration(m.JobHistTab.Stats.MaxWait, m.Log))

	str += "\n"
	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Run times (finished jobs):"))
	str += "\n\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinRun", HumanizeDuration(m.JobHistTab.Stats.MinRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "AvgRun", HumanizeDuration(m.JobHistTab.Stats.AvgRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MedRun", HumanizeDuration(m.JobHistTab.Stats.MedRun, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MaxRun", HumanizeDuration(m.JobHistTab.Stats.MaxRun, m.Log))

	return str
}

func (m Model) JobTabStats() string {

	m.Log.Printf("JobTabStats called\n")

	//str := "Queue statistics (filtered):\n\n"
	str := styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Job states (filtered):"))
	str += "\n\n"

	str += GenCountStrVert(m.JobTab.Stats.StateCnt, m.Log)

	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Pending jobs:"))
	str += "\n\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinWait", HumanizeDuration(m.JobTab.Stats.MinWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "AvgWait", HumanizeDuration(m.JobTab.Stats.AvgWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MedWait", HumanizeDuration(m.JobTab.Stats.MedWait, m.Log))
	str += fmt.Sprintf("%-10s : %s\n", "MaxWait", HumanizeDuration(m.JobTab.Stats.MaxWait, m.Log))

	str += "\n"
	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Running jobs:"))
	str += "\n\n"
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
		// Top Main
		MainWindow.WriteString(fmt.Sprintf("Filter: %10.20s\tItems: %d\n", m.JobTab.Filter.Value(), len(m.JobTab.SqueueFiltered.Jobs)))
		// This is gone to Stats box
		//MainWindow.WriteString(GenCountStr(m.JobTab.Stats.StateCnt, m.Log))
		MainWindow.WriteString("\n")

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

		// Low Main: nil || info || filter || counts
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			// filter
			MainWindow.WriteString("\n")
			MainWindow.WriteString("Filter value (search across: JobID, JobName, Account, UserName, JobState!):\n")
			MainWindow.WriteString(fmt.Sprintf("%s\n", m.JobTab.Filter.View()))
			MainWindow.WriteString("(Enter to apply, Esc to clear filter and abort, Regular expressions supported, syntax details: https://golang.org/s/re2syntax)\n")
		case m.JobTab.InfoOn:
			// info
			MainWindow.WriteString("\n")
			MainWindow.WriteString(styles.JobInfoBox.Render(m.getJobInfo()))
		case m.JobTab.CountsOn:
			// Counts on
			MainWindow.WriteString("\n")
			MainWindow.WriteString(styles.JobInfoBox.Render(m.getJobCounts()))
		}

	case tabJobHist:
		// If sacct timed out/errored, instruct the user to reduce fetch period from default 7 days
		m.Log.Printf("HistFetch: %t HistFetchFail: %t\n", m.JobHistTab.HistFetched, m.JobHistTab.HistFetchFail)
		if m.JobHistTab.HistFetchFail {
			msg := fmt.Sprintf("Fetching jobs history timed out (-t %d seconds), probably too many jobs in the last -d %d days.\n", m.Globals.JobHistTimeout, m.Globals.JobHistStart)
			MainWindow.WriteString(msg)
			MainWindow.WriteString("You can reduce the history period with -d N (days) switch, or increase the history fetch timeout with -t N (seconds) switch.\n")
			break
		}

		// Check if history is here, if not, return "Waiting for sacct..."
		if !m.JobHistTab.HistFetched {
			MainWindow.WriteString("Waiting for job history...\n")
			break
		}

		// Top Main
		MainWindow.WriteString(fmt.Sprintf("Filter: %10.20s\tItems: %d\n", m.JobHistTab.Filter.Value(), len(m.JobHistTab.SacctHistFiltered.Jobs)))
		MainWindow.WriteString("\n")

		// Mid Main: table || table+stats
		switch {
		case m.JobHistTab.StatsOn:
			// table + stats
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, m.tabJobHist(), styles.MenuBoxStyle.Render(m.JobHistTabStats())))
		default:
			// table
			MainWindow.WriteString(m.tabJobHist())
		}

		// Low Main: nil || filter || counts
		switch {
		case m.FilterSwitch == FilterSwitch(m.ActiveTab):
			// filter
			MainWindow.WriteString("\n")
			MainWindow.WriteString("Filter value (search across: JobID, JobName, AccountName, UserName, JobState):\n")
			MainWindow.WriteString(fmt.Sprintf("%s\n", m.JobHistTab.Filter.View()))
			MainWindow.WriteString("(Enter to apply, Esc to clear filter and abort, Regular expressions supported, syntax details: https://golang.org/s/re2syntax)\n")
		case m.JobHistTab.CountsOn:
			// Counts on
			MainWindow.WriteString("\n")
			MainWindow.WriteString(styles.JobInfoBox.Render(m.getJobHistCounts()))
		}

	case tabJobDetails:
		MainWindow.WriteString(m.tabJobDetails())

	case tabJobFromTemplate:
		MainWindow.WriteString(m.tabJobFromTemplate())

	case tabCluster:
		// Top Main
		MainWindow.WriteString(fmt.Sprintf("Filter: %10.20s\tItems: %d\n\n", m.JobClusterTab.Filter.Value(), len(m.JobClusterTab.SinfoFiltered.Nodes)))
		//MainWindow.WriteString(GenCountStr(m.JobClusterTab.Stats.StateCnt, m.Log))
		MainWindow.WriteString(m.tabClusterBars())

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
			MainWindow.WriteString("\n")
			// filter
			MainWindow.WriteString("\n")
			MainWindow.WriteString("Filter value (search across: Name, State, StateFlags!):\n")
			MainWindow.WriteString(fmt.Sprintf("%s\n", m.JobClusterTab.Filter.View()))
			MainWindow.WriteString("(Enter to apply, Esc to clear filter and abort, Regular expressions supported, syntax details: https://golang.org/s/re2syntax)\n")
		default:
			MainWindow.WriteString("\n")
			MainWindow.WriteString(GenCountStr(m.JobClusterTab.Stats.StateCnt, m.Log))
		}

	case tabAbout:
		MainWindow.WriteString(m.tabAbout())
	}

	// FOOTER
	scr.WriteString(lipgloss.JoinVertical(lipgloss.Left, styles.MainWindow.Render(MainWindow.String()), styles.HelpWindow.Render(m.Help.View(keybindings.DefaultKeyMap))))

	return scr.String()
}
