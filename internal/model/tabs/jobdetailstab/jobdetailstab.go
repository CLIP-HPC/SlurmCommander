package jobdetailstab

import (
	"github.com/pja237/slurmcommander-dev/internal/slurm"
)

type JobDetailsTab struct {
	SelJobID string
	slurm.SacctSingleJobHist
}
