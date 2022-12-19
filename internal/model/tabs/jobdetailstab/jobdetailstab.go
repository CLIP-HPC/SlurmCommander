package jobdetailstab

import (
	"github.com/charmbracelet/bubbles/viewport"
	"github.com/CLIP-HPC/SlurmCommander/internal/slurm"
)

type JobDetailsTab struct {
	SelJobID    string
	SelJobIDNew int
	ViewPort    viewport.Model
	slurm.SacctSingleJobHist
}
