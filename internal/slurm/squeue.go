package slurm

import (
	"github.com/CLIP-HPC/SlurmCommander/internal/openapi"
)

type SqueueJSON struct {
	Jobs []openapi.V0039JobResponseProperties
}
