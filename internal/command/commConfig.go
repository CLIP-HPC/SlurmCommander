package command

const (
	tick = 3
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
