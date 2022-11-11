package slurm

import (
	"github.com/pja237/slurmcommander-dev/internal/openapi"
)

type SinfoJSON struct {
	Nodes []openapi.V0039Node
}
