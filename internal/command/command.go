package command

import (
	"encoding/json"
	"log"
	"os/exec"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander/internal/slurm"
)

const (
	tick      = 3
	squeueCmd = "/home/pja/src/go-projects/SlurmCommander/scripts/squeue"
	sinfoCmd  = "/home/pja/src/go-projects/SlurmCommander/scripts/sinfo"
	sacctCmd  = "/home/pja/src/go-projects/SlurmCommander/scripts/sacct-all"
	//sacctJobCmd = "/home/pja/src/go-projects/SlurmCommander/scripts/sacct-46634109"
	sacctJobCmd = "/home/pja/src/go-projects/SlurmCommander/scripts/sacct-46634111"
)

// Calls `sinfo` to get node information for Cluster Tab
func GetSinfo(t time.Time) tea.Msg {
	var siJson slurm.SinfoJSON

	out, err := exec.Command(sinfoCmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Error exec sinfo: %q\n", err)
	}

	err = json.Unmarshal(out, &siJson)
	if err != nil {
		log.Fatalf("Error unmarshall: %q\n", err)
	}

	return siJson

}

func TimedGetSinfo() tea.Cmd {
	// TODO: make timers configurable
	return tea.Tick(tick*time.Second, GetSinfo)
}

func QuickGetSinfo() tea.Cmd {
	// TODO: make timers configurable
	return tea.Tick(0*time.Second, GetSinfo)
}

// Calls `squeue` to get job information for Jobs Tab
func GetSqueue(t time.Time) tea.Msg {

	var sqJson slurm.SqueueJSON

	out, err := exec.Command(squeueCmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Error exec squeue: %q\n", err)
	}

	err = json.Unmarshal(out, &sqJson)
	if err != nil {
		log.Fatalf("Error unmarshall: %q\n", err)
	}

	return sqJson
}

func TimedGetSqueue() tea.Cmd {
	// TODO: make timers configurable
	return tea.Tick(tick*time.Second, GetSqueue)
}

func QuickGetSqueue() tea.Cmd {
	// TODO: make timers configurable
	return tea.Tick(0*time.Second, GetSqueue)
}

// Calls `sacct` to get job information for Jobs History Tab
func GetSacct(t time.Time) tea.Msg {

	var (
		sacct slurm.SacctList
	)

	re, err := regexp.Compile(`\D+`)
	if err != nil {
		// TODO: what do do now? return nil msg?
		return nil
	}

	// run sacct to get list of last 7 days
	out, err := exec.Command(sacctCmd).CombinedOutput()
	if err != nil {
		log.Fatalf("Error exec sacct: %q\n", err)
	}

	// split output into lines
	for _, v := range strings.Split(string(out), "\n") {
		if v == "" {
			continue
		}
		// Test each line if jobid is num-only and append those lines to SacctList
		tmp := strings.Split(v, "|")
		if !re.MatchString(tmp[0]) {
			sacct = append(sacct, tmp)
		}
	}

	return sacct
}

func TimedGetSacct() tea.Cmd {
	return tea.Tick(tick*time.Second, GetSacct)
}

func QuickGetSacct() tea.Cmd {
	return tea.Tick(0*time.Second, GetSacct)
}

func SingleJobGetSacct(jobid string) tea.Cmd {
	return func() tea.Msg {
		var sacctJob slurm.SacctJob

		out, err := exec.Command(sacctJobCmd).CombinedOutput()
		if err != nil {
			log.Fatalf("Error exec sacct: %q\n", err)
		}

		err = json.Unmarshal(out, &sacctJob)
		if err != nil {
			log.Fatalf("Error unmarshall: %q\n", err)
		}

		return sacctJob
	}
}

//func GetJobSacct(jobid string) tea.Msg {
//	var sacctJob slurm.SacctJob
//
//	out, err := exec.Command(sacctJobCmd).CombinedOutput()
//	if err != nil {
//		log.Fatalf("Error exec sacct: %q\n", err)
//	}
//
//	err = json.Unmarshal(out, &sacctJob)
//	if err != nil {
//		log.Fatalf("Error unmarshall: %q\n", err)
//	}
//
//	return sacctJob
//}
