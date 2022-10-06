//go:build prod

package command

const (
	tick        = 3
	squeueCmd   = "/bin/squeue"
	sinfoCmd    = "/bin/sinfo"
	sacctCmd    = "/bin/sacct"
	sacctJobCmd = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-46634109"
	//sacctJobCmd = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-46634111"
	Tag = "prod"
)

var (
	sinfoCmdSwitches  = []string{"-a", "--json"}
	squeueCmdSwitches = []string{"-a", "--json"}
	sacctCmdSwitches  = []string{"-n", "-S", `now-7days`, "-o", `jobid%20,jobname%30,partition,state,exitcode`, "-p"}
)
