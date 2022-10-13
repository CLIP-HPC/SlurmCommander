//go:build prod

package command

const (
	tick           = 3
	squeueCmd      = "/bin/squeue"
	sinfoCmd       = "/bin/sinfo"
	sacctCmd       = "/bin/sacct"
	sacctJobCmd    = "/bin/sacct"
	scancelJobCmd  = "/bin/scancel"
	sholdJobCmd    = "/bin/scontrol"
	srequeueJobCmd = "/bin/scontrol"
	Tag            = "prod"
)

var (
	sinfoCmdSwitches       = []string{"-a", "--json"}
	squeueCmdSwitches      = []string{"-a", "--json"}
	sacctCmdSwitches       = []string{"-n", "-S", `now-7days`, "-o", `jobid%20,jobname%30,partition,state,exitcode`, "-p"}
	sacctJobCmdSwitches    = []string{"-n", "-S", `now-7days`, "--json", "-j"}
	scancelJobCmdSwitches  = []string{}
	sholdJobCmdSwitches    = []string{"hold"}
	srequeueJobCmdSwitches = []string{"requeue"}
)
