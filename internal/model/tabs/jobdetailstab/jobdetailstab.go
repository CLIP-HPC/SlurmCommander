package jobdetailstab

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/pja237/slurmcommander-dev/internal/slurm"
)

type JobDetailsTab struct {
	SelJobID    string
	SelJobIDNew int
	ViewPort    viewport.Model
	slurm.SacctSingleJobHist
}
