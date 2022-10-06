//go:build !local && !prod

package command

const (
	tick      = 3
	squeueCmd = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/squeue"
	sinfoCmd  = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sinfo"
	sacctCmd  = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-all"
	//sacctJobCmd = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-46634109"
	sacctJobCmd = "/home/pja/src/go-projects/SlurmCommander-dev/scripts/sacct-46634111"
	Tag         = "!local && !prod"
)
