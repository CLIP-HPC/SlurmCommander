package jobhisttab

import (
	"fmt"
	"log"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/CLIP-HPC/SlurmCommander/internal/generic"
	"github.com/CLIP-HPC/SlurmCommander/internal/styles"
)

func (jh *JobHistTab) tabJobHist() string {

	return jh.SacctTable.View() + "\n"
}

func (jh *JobHistTab) JobHistTabStats(l *log.Logger) string {

	l.Printf("JobHistTabStats called\n")

	str := styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Historical job states (filtered):"))
	str += "\n"
	str += generic.GenCountStrVert(jh.Stats.StateCnt, l)

	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Waiting times (finished jobs):"))
	str += "\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinWait", generic.HumanizeDuration(jh.Stats.MinWait, l))
	str += fmt.Sprintf("%-10s : %s\n", "AvgWait", generic.HumanizeDuration(jh.Stats.AvgWait, l))
	str += fmt.Sprintf("%-10s : %s\n", "MedWait", generic.HumanizeDuration(jh.Stats.MedWait, l))
	str += fmt.Sprintf("%-10s : %s\n", "MaxWait", generic.HumanizeDuration(jh.Stats.MaxWait, l))

	str += "\n"
	str += styles.StatsSeparatorTitle.Render(fmt.Sprintf("%-30s", "Run times (finished jobs):"))
	str += "\n"
	str += fmt.Sprintf("%-10s : %s\n", " ", "dd-hh:mm:ss")
	str += fmt.Sprintf("%-10s : %s\n", "MinRun", generic.HumanizeDuration(jh.Stats.MinRun, l))
	str += fmt.Sprintf("%-10s : %s\n", "AvgRun", generic.HumanizeDuration(jh.Stats.AvgRun, l))
	str += fmt.Sprintf("%-10s : %s\n", "MedRun", generic.HumanizeDuration(jh.Stats.MedRun, l))
	str += fmt.Sprintf("%-10s : %s\n", "MaxRun", generic.HumanizeDuration(jh.Stats.MaxRun, l))

	return str
}

func (jh *JobHistTab) getJobHistCounts() string {
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
	for _, v := range jh.Breakdowns.Top5user {
		top5u += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	top5a += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Top 5 Accounts"))
	top5a += "\n"
	for _, v := range jh.Breakdowns.Top5acc {
		top5a += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	jpp += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Jobs per Partition"))
	jpp += "\n"
	for _, v := range jh.Breakdowns.JobPerPart {
		jpp += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	jpq += styles.TextYellowOnBlue.Render(fmt.Sprintf(fmtTitle, "Jobs per QoS"))
	jpq += "\n"
	for _, v := range jh.Breakdowns.JobPerQos {
		jpq += fmt.Sprintf(fmtStr, v.Name, v.Count)
	}

	top5u = styles.CountsBox.Render(top5u)
	top5a = styles.CountsBox.Render(top5a)
	jpq = styles.CountsBox.Render(jpq)
	jpp = styles.CountsBox.Render(jpp)

	ret = lipgloss.JoinHorizontal(lipgloss.Top, top5u, top5a, jpp, jpq)

	return ret
}

func (jh *JobHistTab) View(l *log.Logger) string {
	var (
		Header     strings.Builder
		MainWindow strings.Builder
	)

	// Show parameters to the user
	heading := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FAFAFA")).Background(lipgloss.Color("#37ABBF")).Width(12).PaddingLeft(1).PaddingRight(1).Align(lipgloss.Right)
	Header.WriteString(fmt.Sprintf("%s Start: %10.20s\tEnd: %10.20s\tTimeout: %d\n", heading.Render("Parameters"), jh.JobHistStart, jh.JobHistEnd, jh.JobHistTimeout))

	// If sacct timed out/errored, instruct the user to reduce fetch period from default 7 days
	l.Printf("HistFetch: %t HistFetchFail: %t\n", jh.HistFetched, jh.HistFetchFail)
	if jh.HistFetchFail {
		Header.WriteString(fmt.Sprintf("Fetching jobs history failed! Maybe the timed out was too short(%d seconds)? See the returned error message:\n", jh.JobHistTimeout))
		Header.WriteString("You can you can modify the time ranges or timeout using the 'time-ranges' menu\n")
	}

	// Check if history is here, if not, return "Waiting for sacct..."
	if !jh.HistFetchFail && !jh.HistFetched {
		Header.WriteString("Waiting for job history...\n")
		return Header.String()
	}

	if !jh.HistFetchFail {
		// Rest of header
		Header.WriteString(fmt.Sprintf("%s Query: %10.20s\tItems: %d\n", heading.Render("Filter"), jh.Filter.Value(), len(jh.SacctHistFiltered.Jobs)))
		Header.WriteString("\n")

		// Top Main
		MainWindow.WriteString(jh.tabJobHist())
	} else {
		MainWindow.WriteString("\n")
	}

	// Next we join table Vertically with: nil || filter || params || counts
	switch {
	case jh.FilterOn:
		MainWindow.WriteString("\n")
		MainWindow.WriteString("Filter value (search in joined: JobID + JobName + QoS + AccountName + UserName + JobState):\n")
		MainWindow.WriteString(fmt.Sprintf("%s\n", jh.Filter.View()))
		MainWindow.WriteString("(Enter to apply, Esc to clear filter and abort, Regular expressions supported, syntax details: https://golang.org/s/re2syntax)\n")

	case jh.UserInputsOn:
		MainWindow.WriteString("\n")
		MainWindow.WriteString(fmt.Sprintf("Command Parameters:\n"))
		for i := range jh.UserInputs.Params {
			MainWindow.WriteString(fmt.Sprintf("%s: %s\n", jh.UserInputs.ParamTexts[i], jh.UserInputs.Params[i].View()))
		}
		MainWindow.WriteString("(Enter to apply, or Esc to clear params and abort)\n")

	case jh.CountsOn:
		// Counts on
		MainWindow.WriteString("\n")
		MainWindow.WriteString(styles.JobInfoBox.Render(jh.getJobHistCounts()))
	}

	// Last, if needed we join Stats Horizontally with Main
	switch {
	case jh.StatsOn:
		// table + stats
		X := MainWindow.String()
		MainWindow.Reset()
		MainWindow.WriteString(lipgloss.JoinHorizontal(lipgloss.Top, X, styles.StatsBoxStyle.Render(jh.JobHistTabStats(l))))
	}

	return Header.String() + MainWindow.String()
}
