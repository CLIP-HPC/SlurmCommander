package command

import (
	"bufio"
	"encoding/json"
	"log"
	"os/exec"
	"os/user"
	"regexp"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/config"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
)

var cc config.ConfigContainer

// NewCmdCC sets the package variable ConfigContainer to the values from config files/defaults.
// To be used in the package locally without being passed around.
// Not the smartest choice, consider refactoring later.
func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

// UserName is a linux username string.
type UserName string

// GetUserName returns the linux username string.
// Used to call sacctmgr and fetch users associations.
func GetUserName(l *log.Logger) tea.Cmd {

	return func() tea.Msg {

		l.Printf("Fetching UserName\n")
		u, err := user.Current()
		if err != nil {
			l.Fatalf("GetUserName FAILED: %s", err)
		}

		l.Printf("Return UserName: %s\n", u.Username)
		return UserName(u.Username)
	}
}

// UserAssoc is a list of associations between current username and slurm accounts.
type UserAssoc []string

// GetUserAssoc returns a list of associations between current username and slurm accounts.
// Used later in sacct --json call to fetch users job history
func GetUserAssoc(u string, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var (
			ua UserAssoc
			sw []string
		)

		cmd := cc.Binpaths["sacctmgr"]
		sw = append(sw, sacctmgrCmdSwitches...)
		sw = append(sw, "user="+u)
		//c := exec.Command(sacctmgrCmd, sw...)
		c := exec.Command(cmd, sw...)

		l.Printf("GetUserAssoc about to run: %v %v\n", cmd, sw)
		stdOut, err := c.StdoutPipe()
		if err != nil {
			l.Fatalf("StdoutPipe call FAILED with %s\n", err)
		}
		if e := c.Start(); e != nil {
			l.Fatalf("cmd.Run call FAILED with %s\n", err)
		}
		s := bufio.NewScanner(stdOut)
		for s.Scan() {
			l.Printf("Got UserAssoc %s -> %s\n", u, s.Text())
			ua = append(ua, s.Text())
		}
		if e := c.Wait(); e != nil {
			l.Fatalf("cmd.Wait call FAILED with %s\n", err)

		}

		return ua
	}
}

// Calls `sinfo` to get node information for Cluster Tab
func GetSinfo(t time.Time) tea.Msg {
	var siJson slurm.SinfoJSON

	cmd := cc.Binpaths["sinfo"]
	out, err := exec.Command(cmd, sinfoCmdSwitches...).CombinedOutput()
	if err != nil {
		log.Fatalf("Error exec sinfo: %s : %q\n", cmd, err)
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

	cmd := cc.Binpaths["squeue"]
	out, err := exec.Command(cmd, squeueCmdSwitches...).CombinedOutput()
	if err != nil {
		log.Fatalf("Error exec squeue: %s : %q\n", cmd, err)
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
	cmd := cc.Binpaths["sacct"]
	out, err := exec.Command(cmd, sacctCmdSwitches...).CombinedOutput()
	if err != nil {
		log.Fatalf("Error exec sacct: %s : %q\n", cmd, err)
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

func GetSacctHist(uaccs string, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var (
			sacctJobs slurm.SacctJobHist
			sw        []string
		)

		// TODO: use https://pkg.go.dev/context#WithTimeout
		// ...to limit the sacct runs. Flip bool switch in model to failed if it exceeds.
		// Setup command line switch to pick how many days of sacct to fetch in case of massive runs.

		l.Printf("GetSacctHist(%q) start\n", uaccs)
		cmd := cc.Binpaths["sacct"]
		sw = append(sacctHistCmdSwitches, "-A", uaccs)
		l.Printf("EXEC: %q %q\n", cmd, sw)
		out, err := exec.Command(cmd, sw...).CombinedOutput()
		l.Printf("EXEC returned: %d bytes\n", len(out))
		if err != nil {
			l.Fatalf("Error exec sacct: %q\n", err)
		}

		err = json.Unmarshal(out, &sacctJobs)
		if err != nil {
			l.Fatalf("Error unmarshall: %q\n", err)
		}

		l.Printf("Unmarshalled %d jobs from hist\n", len(sacctJobs.Jobs))

		return sacctJobs
	}
}

func SingleJobGetSacct(jobid string, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var sacctJob slurm.SacctSingleJobHist

		l.Printf("SingleJobGetSacct start.\n")
		cmd := cc.Binpaths["sacct"]
		switches := append(sacctJobCmdSwitches, jobid)
		out, err := exec.Command(cmd, switches...).CombinedOutput()
		l.Printf("SingleJobGetSacct EXEC: %q %q\n", cmd, switches)
		if err != nil {
			log.Fatalf("Error exec sacct: %q\n", err)
		}
		l.Printf("SingleJobGetSacct EXEC got bytes: %d\n", len(out))

		err = json.Unmarshal(out, &sacctJob)
		if err != nil {
			log.Fatalf("Error unmarshall: %q\n", err)
		}
		l.Printf("SingleJobGetSacct UNMARSHALLED items: %d\n", len(sacctJob.Jobs))

		// do a new type for a single job, otherwise the update() will handle single job and the whole hist list in same code-path
		return sacctJob
	}
}

type ScancelSent struct {
	Jobid string
}

func CallScancel(jobid string, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var scret ScancelSent = ScancelSent{
			Jobid: jobid,
		}

		cmd := cc.Binpaths["scancel"]
		switches := append(scancelJobCmdSwitches, jobid)

		l.Printf("EXEC: %q %q\n", cmd, switches)
		out, err := exec.Command(cmd, switches...).CombinedOutput()
		if err != nil {
			l.Fatalf("Error exec scancel: %q\n", err)
		}
		l.Printf("EXEC output: %q\n", out)

		return scret
	}
}

type SHoldSent struct {
	Jobid string
}

func CallScontrolHold(jobid string, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var scret SHoldSent = SHoldSent{
			Jobid: jobid,
		}

		cmd := cc.Binpaths["scontrol"]
		switches := append(sholdJobCmdSwitches, jobid)

		l.Printf("EXEC: %q %q\n", cmd, switches)
		out, err := exec.Command(cmd, switches...).CombinedOutput()
		if err != nil {
			l.Fatalf("Error exec hold: %q\n", err)
		}
		l.Printf("EXEC output: %q\n", out)

		return scret
	}
}

type SRequeueSent struct {
	Jobid string
}

// TODO: unify this to a single function
func CallScontrolRequeue(jobid string, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var scret SRequeueSent = SRequeueSent{
			Jobid: jobid,
		}

		cmd := cc.Binpaths["scontrol"]
		switches := append(srequeueJobCmdSwitches, jobid)

		l.Printf("EXEC: %q %q\n", cmd, switches)
		out, err := exec.Command(cmd, switches...).CombinedOutput()
		if err != nil {
			//l.Fatalf("Error exec requeue: %q\n", err)
			l.Printf("Error exec requeue: %q\n", err)
			l.Printf("Possible Reason: requeue can only be executed for batch jobs (e.g. won't work on srun --pty)\n")
		}
		l.Printf("EXEC output: %q\n", out)

		return scret
	}
}

type SBatchSent struct {
	JobFile string
}

// TODO: unify this to a single function
func CallSbatch(jobfile string, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var scret SBatchSent = SBatchSent{
			JobFile: jobfile,
		}

		cmd := cc.Binpaths["sbatch"]
		switches := append(sbatchCmdSwitches, jobfile)

		l.Printf("EXEC: %q %q\n", cmd, switches)
		out, err := exec.Command(cmd, switches...).CombinedOutput()
		if err != nil {
			l.Fatalf("Error exec sbatch: %q\n", err)
		}
		l.Printf("EXEC output: %q\n", out)

		return scret
	}
}
