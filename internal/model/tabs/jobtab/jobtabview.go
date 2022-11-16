package jobtab

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/pja237/slurmcommander-dev/internal/generic"
	"github.com/pja237/slurmcommander-dev/internal/styles"
)

func (jt *JobTab) tabJobs() string {

	return jt.SqueueTable.View() + "\n"
}

func (jt *JobTab) getJobInfo(l *log.Logger) string {
	var scr strings.Builder

	n := jt.SqueueTable.Cursor()
	l.Printf("getJobInfo: cursor at %d table rows: %d\n", n, len(jt.SqueueFiltered.Jobs))
	if len(jt.SqueueFiltered.Jobs) == 0 || n == -1 {
		return "Select a job"
	}

	fmtStr := "%-15s : %-30s\n"
	fmtStrLast := "%-15s : %-30s"
	//ibFmt := "Job Name: %s\nJob Command: %s\nOutput: %s\nError: %s\n"
	infoBoxLeft := fmt.Sprintf(fmtStr, "Partition", *jt.SqueueFiltered.Jobs[n].Partition)
	infoBoxLeft += fmt.Sprintf(fmtStr, "QoS", *jt.SqueueFiltered.Jobs[n].Qos)
	infoBoxLeft += fmt.Sprintf(fmtStr, "TRES", *jt.SqueueFiltered.Jobs[n].TresReqStr)
	infoBoxLeft += fmt.Sprintf(fmtStr, "Batch Host", *jt.SqueueFiltered.Jobs[n].BatchHost)
	if jt.SqueueFiltered.Jobs[n].JobResources.Nodes != nil {
		infoBoxLeft += fmt.Sprintf(fmtStrLast, "AllocNodes", *jt.SqueueFiltered.Jobs[n].JobResources.Nodes)
	} else {
		infoBoxLeft += fmt.Sprintf(fmtStrLast, "AllocNodes", "none")

	}

	infoBoxRight := fmt.Sprintf(fmtStr, "Array Job ID", strconv.Itoa(*jt.SqueueFiltered.Jobs[n].ArrayJobId))
	if jt.SqueueFiltered.Jobs[n].ArrayTaskId != nil {
		infoBoxRight += fmt.Sprintf(fmtStr, "Array Task ID", strconv.Itoa(*jt.SqueueFiltered.Jobs[n].ArrayTaskId))
	} else {
		infoBoxRight += fmt.Sprintf(fmtStr, "Array Task ID", "NoTaskID")
	}
	infoBoxRight += fmt.Sprintf(fmtStr, "Gres Details", strings.Join(*jt.SqueueFiltered.Jobs[n].GresDetail, ","))
	infoBoxRight += fmt.Sprintf(fmtStr, "Features", *jt.SqueueFiltered.Jobs[n].Features)
	infoBoxRight += fmt.Sprintf(fmtStrLast, "wckey", *jt.SqueueFiltered.Jobs[n].Wckey)

	infoBoxMiddle := fmt.Sprintf(fmtStr, "Submit", time.Unix(*jt.SqueueFiltered.Jobs[n].SubmitTime, 0))
	if *jt.SqueueFiltered.Jobs[n].StartTime != 0 {
		infoBoxMiddle += fmt.Sprintf(fmtStr, "Start /expected", time.Unix(*jt.SqueueFiltered.Jobs[n].StartTime, 0))
	} else {
		infoBoxMiddle += fmt.Sprintf(fmtStr, "Start", "unknown")
	}
	// placeholder lines
	infoBoxMiddle += "\n"
	infoBoxMiddle += "\n"
	// EO placeholder lines
	infoBoxMiddle += fmt.Sprintf(fmtStrLast, "State reason", *jt.SqueueFiltered.Jobs[n].StateReason)

	infoBoxWide := fmt.Sprintf(fmtStr, "Job Name", *jt.SqueueFiltered.Jobs[n].Name)
	infoBoxWide += fmt.Sprintf(fmtStr, "Command", *jt.SqueueFiltered.Jobs[n].Command)
	infoBoxWide += fmt.Sprintf(fmtStr, "StdOut", *jt.SqueueFiltered.Jobs[n].StandardOutput)
	infoBoxWide += fmt.Sprintf(fmtStr, "StdErr", *jt.SqueueFiltered.Jobs[n].StandardError)
	infoBoxWide += fmt.Sprintf(fmtStrLast, "Working Dir", *jt.SqueueFiltered.Jobs[n].CurrentWorkingDirectory)

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

func (jt *JobTab) getJobCounts() string {
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
	for _, v := range jt.Breakdowns.Top5user {
		top5u += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	top5a += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Top 5 Accounts"))
	top5a += "\n"
	for _, v := range jt.Breakdowns.Top5acc {
		top5a += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	jpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Jobs per Partition"))
	jpp += "\n"
	for _, v := range jt.Breakdowns.JobPerPart {
		jpp += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	jpq += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Jobs per QoS"))
	jpq += "\n"
	for _, v := range jt.Breakdowns.JobPerQos {
		jpq += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	top5u = styles.CountsBox.Render(top5u)
	top5a = styles.CountsBox.Render(top5a)
	jpq = styles.CountsBox.Render(jpq)
	jpp = styles.CountsBox.Render(jpp)

	ret = lipgloss.JoinHorizontal(lipgloss.Top, top5u, top5a, jpp, jpq)

	return ret
}

func (jt *JobTab) JobTabStats(l *log.Logger) string {

	l.Printf("JobTabStats called\n")

	str := styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Job states (filtered):"))
	str += "\n"

	str += generic.GenCountStrVert(jt.Stats.StateCnt, l)

	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Pending jobs:"))
	str += "\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinWait", generic.HumanizeDuration(jt.Stats.MinWait, l))
	str += fmt.Sprintf("%-10s : %s\n", "AvgWait", generic.HumanizeDuration(jt.Stats.AvgWait, l))
	str += fmt.Sprintf("%-10s : %s\n", "MedWait", generic.HumanizeDuration(jt.Stats.MedWait, l))
	str += fmt.Sprintf("%-10s : %s\n", "MaxWait", generic.HumanizeDuration(jt.Stats.MaxWait, l))

	str += "\n"
	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Running jobs:"))
	str += "\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinRun", generic.HumanizeDuration(jt.Stats.MinRun, l))
	str += fmt.Sprintf("%-10s : %s\n", "AvgRun", generic.HumanizeDuration(jt.Stats.AvgRun, l))
	str += fmt.Sprintf("%-10s : %s\n", "MedRun", generic.HumanizeDuration(jt.Stats.MedRun, l))
	str += fmt.Sprintf("%-10s : %s\n", "MaxRun", generic.HumanizeDuration(jt.Stats.MaxRun, l))

	return str
}

func (jt *JobTab) View(l *log.Logger) string {

	var (
		Header     strings.Builder
		MainWindow strings.Builder
	)

	l.Printf("IN JobTab.View()")

	// Header
	//Header.WriteString("\n")
	Header.WriteString(fmt.Sprintf("Filter: %10.30s\tItems: %d\n", jt.Filter.Value(), len(jt.SqueueFiltered.Jobs)))
	Header.WriteString("\n")

	// Mid Main: table || table+stats || table+menu

	// Case Info OFF
	if !jt.InfoOn {
		// Table always here
		MainWindow.WriteString(jt.tabJobs())

		// Bellow table stich Filter || Counts
		switch {
		case jt.FilterOn:
			// filter
			MainWindow.WriteString("\n")
			MainWindow.WriteString("Filter value (Search in joined: JobID + JobName + Account + UserName + JobState):\n")
			MainWindow.WriteString(fmt.Sprintf("%s\n", jt.Filter.View()))
			MainWindow.WriteString("(Enter to apply, Esc to clear filter and abort, Regular expressions supported, syntax details: https://golang.org/s/re2syntax)\n")
		case jt.CountsOn:
			// Counts on
			MainWindow.WriteString("\n")
			MainWindow.WriteString(styles.JobInfoBox.Render(jt.getJobCounts()))
		}

		// Then join all that Horizontally with Menu || Stats
		switch {
		case jt.MenuOn:
			X := MainWindow.String()
			MainWindow.Reset()
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, X, styles.MenuBoxStyle.Render(jt.Menu.View())))
			l.Printf("\nITEMS LIST: %#v\n", jt.Menu.Items())
		case jt.StatsOn:
			X := MainWindow.String()
			MainWindow.Reset()
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, X, styles.MenuBoxStyle.Render(jt.JobTabStats(l))))
		}
	} else {
		// Case Info ON

		// First join Horizontally Table with Menu || Stats
		switch {
		case jt.MenuOn:
			// table + menu
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, jt.tabJobs(), styles.MenuBoxStyle.Render(jt.Menu.View())))
			l.Printf("\nITEMS LIST: %#v\n", jt.Menu.Items())
		case jt.StatsOn:
			// table + stats
			MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, jt.tabJobs(), styles.MenuBoxStyle.Render(jt.JobTabStats(l))))
		default:
			// table
			MainWindow.WriteString(jt.tabJobs())
		}

		// Then stich the Filter || Info || Counts below
		switch {
		case jt.FilterOn:
			// filter
			MainWindow.WriteString("\n")
			MainWindow.WriteString("Filter value (Search in joined: JobID + JobName + Account + UserName + JobState):\n")
			MainWindow.WriteString(fmt.Sprintf("%s\n", jt.Filter.View()))
			MainWindow.WriteString("(Enter to apply, Esc to clear filter and abort, Regular expressions supported, syntax details: https://golang.org/s/re2syntax)\n")
		case jt.InfoOn:
			// info
			MainWindow.WriteString("\n")
			MainWindow.WriteString(styles.JobInfoBox.Render(jt.getJobInfo(l)))
		case jt.CountsOn:
			// Counts on
			MainWindow.WriteString("\n")
			MainWindow.WriteString(styles.JobInfoBox.Render(jt.getJobCounts()))
		}
	}

	return Header.String() + MainWindow.String()
}
