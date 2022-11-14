package command

import (
	"bufio"
	"fmt"
	"log"
	"os/exec"
	"os/user"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/pja237/slurmcommander-dev/internal/config"
)

var cc config.ConfigContainer

// NewCmdCC sets the package variable ConfigContainer to the values from config files/defaults.
// To be used in the package locally without being passed around.
// Not the smartest choice, consider refactoring later.
func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

type ErrorMsg struct {
	From    string
	ErrHelp string
	OrigErr error
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
			//l.Fatalf("GetUserName FAILED: %s", err)
			return ErrorMsg{
				From:    "GetUserName",
				ErrHelp: "Failed to get username, hard to imagine why. Please open an issue with us here: https://github.com/pja237/SlurmCommander-dev/issues/new/choose",
				OrigErr: err,
			}
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
			//l.Fatalf("StdoutPipe call FAILED with %s\n", err)
			return ErrorMsg{
				From:    "GetUserAssoc",
				ErrHelp: "Failed setting up command StdoutPipe()",
				OrigErr: err,
			}
		}

		if e := c.Start(); e != nil {
			//l.Fatalf("cmd.Run call FAILED with %s\n", err)
			return ErrorMsg{
				From:    "GetUserAssoc",
				ErrHelp: "Failed Start()ing sacctmgr",
				OrigErr: err,
			}
		}

		s := bufio.NewScanner(stdOut)
		for s.Scan() {
			l.Printf("Got UserAssoc %s -> %s\n", u, s.Text())
			ua = append(ua, s.Text())
		}
		if e := c.Wait(); e != nil {
			//l.Fatalf("cmd.Wait call FAILED with %s\n", err)
			return ErrorMsg{
				From:    "GetUserAssoc",
				ErrHelp: "Failed Wait()ing for sacctmgr to exit",
				OrigErr: err,
			}
		}

		return ua
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
			//l.Fatalf("Error exec scancel: %q\n", err)
			return ErrorMsg{
				From:    "CallScancel",
				ErrHelp: "Failed to run scancel: check command paths in scom.conf, check that you have permissions to cancel it",
				OrigErr: err,
			}
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
			//l.Fatalf("Error exec hold: %q\n", err)
			return ErrorMsg{
				From:    "CallScontrolHold",
				ErrHelp: "Failed to run scontrol hold : check command paths in scom.conf, check that you have permissions to hold it",
				OrigErr: err,
			}
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
			return ErrorMsg{
				From:    "CallScontrolRequeue",
				ErrHelp: "Failed to run scontrol requeue: check command paths in scom.conf, check that you have permissions to requeue it, check that it's not srun --pty type",
				OrigErr: err,
			}
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
			return ErrorMsg{
				From:    "CallSbatch",
				ErrHelp: "Failed to run sbatch, check scom.conf and the paths there.",
				OrigErr: err,
			}
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
		return ErrorMsg{
			From:    "CallSsh",
			ErrHelp: fmt.Sprintf("Failed ssh to %s, possible reasons: you don't have a job on the node, ssh not allowed at all, etc.", node),
			OrigErr: err,
		}
		//return SshCompleted{
		//	SshErr: err,
		//}
	})
}
