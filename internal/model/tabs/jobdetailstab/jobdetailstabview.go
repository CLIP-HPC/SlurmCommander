package jobdetailstab

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/pja237/slurmcommander-dev/internal/model/tabs/jobhisttab"
	"github.com/pja237/slurmcommander-dev/internal/styles"
)

func (jd *JobDetailsTab) tabJobDetails(jh *jobhisttab.JobHistTab, l *log.Logger) (scr string) {

	var (
		runT  time.Duration
		waitT time.Duration
	)

	// race between View() call and command.SingleJobGetSacct(jd.SelJobID) call
	switch {
	//case jd.SelJobID == "":
	//	return "Select a job from the Job History tab.\n"
	case jd.SelJobIDNew == -1:
		return "Select a job from the Job History tab.\n"
	//case len(m.SacctSingleJobHist.Jobs) == 0:
	//	return fmt.Sprintf("Waiting for job %s info...\n", jd.SelJobID)
	case len(jh.SacctHistFiltered.Jobs) == 0:
		//return fmt.Sprintf("Waiting for job %s info...\n", jd.SelJobID)
		return "Select a job from the Job History tab.\n"
	}

	//width := m.Globals.winW - 10

	//job := m.SacctSingleJobHist.Jobs[0]
	// NEW:
	job := jh.SacctHistFiltered.Jobs[jd.SelJobIDNew]

	l.Printf("Job Details req %#v ,got: %#v\n", jd.SelJobID, job.JobId)

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

		l.Printf("Job Details, step: %d name: %s\n", i, *v.Step.Name)
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
		l.Printf("Dumping step allocation: %#v\n", *v.Tres.Allocated)
		l.Printf("ALLOCATED:\n")
		for i, t := range *v.Tres.Allocated {
			if t.Count != nil {
				l.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
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
		//l.Printf("REQ:Min\n")
		//for i, t := range *v.Tres.Requested.Min {
		//	if t.Count != nil {
		//		l.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresReqMin += " "
		//		tresReqMin += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// REQUESTED:MAX
		//l.Printf("REQ:Max\n")
		//tresReqMax += "REQUESTED:Max:\n"
		//for i, t := range *v.Tres.Requested.Min {
		//	if t.Count != nil {
		//		l.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresReqMax += " "
		//		tresReqMax += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// REQUESTED:AVG
		//l.Printf("REQ:Avg\n")
		//tresReqAvg += "REQUESTED:Avg:\n"
		//for i, t := range *v.Tres.Requested.Average {
		//	if t.Count != nil {
		//		l.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresReqAvg += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// REQUESTED:TOT
		//tresReqAvg += "REQUESTED:Tot:\n"
		//l.Printf("REQ:Tot\n")
		//for i, t := range *v.Tres.Requested.Total {
		//	if t.Count != nil {
		//		l.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresReqTotal += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// Consumed:Min
		//tresConMin += "CONSUMED:Min:\n"
		//l.Printf("CONS:Min\n")
		//for i, t := range *v.Tres.Consumed.Min {
		//	if t.Count != nil {
		//		l.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
		//		tresConMin += fmt.Sprintf(fmtStr, *t.Type, strconv.Itoa(*t.Count))
		//	}
		//}
		//// Consumed:Max
		//tresConMax += "CONSUMED:Max:\n"
		//l.Printf("CONS:Max\n")
		//for i, t := range *v.Tres.Consumed.Max {
		//	if t.Count != nil {
		//		l.Printf("Dumping type %d : %s - %d\n", i, *t.Type, *t.Count)
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

func (jd *JobDetailsTab) SetViewportContent(jh *jobhisttab.JobHistTab, l *log.Logger) string {
	var (
		MainWindow strings.Builder
	)

	jd.ViewPort.SetContent(jd.tabJobDetails(jh, l))

	return MainWindow.String()
}

func (jd *JobDetailsTab) View(jh *jobhisttab.JobHistTab, l *log.Logger) string {
	var (
		MainWindow strings.Builder
	)

	MainWindow.WriteString(jd.ViewPort.View())

	return MainWindow.String()
}
