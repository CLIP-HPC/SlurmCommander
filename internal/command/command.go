package command

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"os/user"

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
		sw = append(sw, SacctmgrCmdSwitches...)
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

func SingleJobGetSacct(jobid string, d uint, l *log.Logger) tea.Cmd {
	// TODO: this we DO NOT need!
	return func() tea.Msg {
		var sacctJob slurm.SacctSingleJobHist

		l.Printf("SingleJobGetSacct start.\n")
		cmd := cc.Binpaths["sacct"]
		start := fmt.Sprintf("now-%ddays", d)
		switches := append(SacctJobCmdSwitches, jobid)
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
		switches := append(ScancelJobCmdSwitches, jobid)

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
		switches := append(SholdJobCmdSwitches, jobid)

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
		switches := append(SrequeueJobCmdSwitches, jobid)

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
		switches := append(SbatchCmdSwitches, jobfile)

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
