//go:build local

package command

const (
	tick      = 3
	squeueCmd = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/squeue"
	sinfoCmd  = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sinfo"
	sacctCmd  = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-all"
	//sacctJobCmd = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-46634109"
	sacctJobCmd = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-46634111"
	Tag         = "local"
)

var (
	sinfoCmdSwitches  = []string{"-a", "--json"}
	squeueCmdSwitches = []string{"-a", "--json"}
	sacctCmdSwitches  = []string{"-n", "-S", `now-7days`, "-o", `jobid%20,jobname%30,partition,state,exitcode`, "-p"}
)
