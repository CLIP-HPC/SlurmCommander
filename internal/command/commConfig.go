package command

const (
	tick = 3
)

var (
	sinfoCmdSwitches       = []string{"-a", "--json"}
	squeueCmdSwitches      = []string{"-a", "--json"}
	sacctHistCmdSwitches   = []string{"-n", "--json"}
	sacctJobCmdSwitches    = []string{"-n", "--json", "-j"}
	scancelJobCmdSwitches  = []string{}
	sholdJobCmdSwitches    = []string{"hold"}
	srequeueJobCmdSwitches = []string{"requeue"}
	sbatchCmdSwitches      = []string{}
	sacctmgrCmdSwitches    = []string{"list", "Association", "format=account", "-P", "-n"}
)
