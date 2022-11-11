package slurm

import (
	"github.com/pja237/slurmcommander-dev/internal/openapi"
)

type SqueueJSON struct {
	Jobs []openapi.V0039JobResponseProperties
}
