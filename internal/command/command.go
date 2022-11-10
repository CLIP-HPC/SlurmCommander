package command

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"os/user"
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

type JobHistTabMsg struct {
	HistFetchFail bool
	slurm.SacctJobHist
}

func GetSacctHist(uaccs string, d uint, t uint, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var (
			jht JobHistTabMsg
		)

		// TODO: Setup command line switch to pick how many days of sacct to fetch in case of massive runs.

		l.Printf("GetSacctHist(%q) start: days %d, timeout: %d\n", uaccs, d, t)

		// setup context with 5 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
		defer cancel()

		// prepare command
		cmd := cc.Binpaths["sacct"]
		start := fmt.Sprintf("now-%ddays", d)
		switches := append(sacctHistCmdSwitches, "-S", start)
		switches = append(switches, "-A", uaccs)

		l.Printf("EXEC: %q %q\n", cmd, switches)
		out, err := exec.CommandContext(ctx, cmd, switches...).CombinedOutput()
		if err != nil {
			l.Printf("Error exec sacct: %q\n", err)
			// set error, return.
			jht.HistFetchFail = true
			return jht
		}
		l.Printf("EXEC returned: %d bytes\n", len(out))

		err = json.Unmarshal(out, &jht.SacctJobHist)
		if err != nil {
			jht.HistFetchFail = true
			l.Printf("Error unmarshall: %q\n", err)
			return jht
		}

		jht.HistFetchFail = false
		l.Printf("Unmarshalled %d jobs from hist\n", len(jht.SacctJobHist.Jobs))

		return jht
	}
}

func SingleJobGetSacct(jobid string, d uint, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var sacctJob slurm.SacctSingleJobHist

		l.Printf("SingleJobGetSacct start.\n")
		cmd := cc.Binpaths["sacct"]
		start := fmt.Sprintf("now-%ddays", d)
		switches := append(sacctJobCmdSwitches, jobid)
		switches = append(switches, "-S", start)

		l.Printf("SingleJobGetSacct EXEC: %q %q\n", cmd, switches)
		out, err := exec.Command(cmd, switches...).CombinedOutput()
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

type SshCompleted struct {
	SshErr error
}

func CallSsh(node string, l *log.Logger) tea.Cmd {
	l.Printf("Start ssh to %s\n", node)
	ssh := exec.Command("ssh", node)
	return tea.ExecProcess(ssh, func(err error) tea.Msg {
		l.Printf("End ssh with error: %s\n", err)
		return SshCompleted{
			SshErr: err,
		}
	})
}
