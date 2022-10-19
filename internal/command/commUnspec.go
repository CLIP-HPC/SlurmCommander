//go:build !local && !prod

package command

const (
	tick           = 3
	squeueCmd      = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/squeue"
	sinfoCmd       = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sinfo"
	sacctCmd       = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-all"
	sacctHistCmd   = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-all"
	sacctJobCmd    = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-46634109"
	scancelJobCmd  = "/bin/scancel"
	sholdJobCmd    = "/bin/scontrol"
	srequeueJobCmd = "/bin/scontrol"
	sbatchCmd      = "/bin/sbatch"
	sacctmgrCmd    = "/bin/sacctmgr"
	Tag            = "!local && !prod"
)

var (
	sinfoCmdSwitches       = []string{"-a", "--json"}
	squeueCmdSwitches      = []string{"-a", "--json"}
	sacctCmdSwitches       = []string{"-n", "-S", `now-7days`, "-o", `jobid%20,jobname%30,partition,state,exitcode`, "-p"}
	sacctHistCmdSwitches   = []string{"-n", "-S", `now-7days`, "--json"}
	sacctJobCmdSwitches    = []string{"-n", "-S", `now-7days`, "--json", "-j"}
	scancelJobCmdSwitches  = []string{}
	sholdJobCmdSwitches    = []string{"hold"}
	srequeueJobCmdSwitches = []string{"requeue"}
	sbatchCmdSwitches      = []string{}
	sacctmgrCmdSwitches    = []string{"list", "Association", "format=account", "-P", "-n"}
)
