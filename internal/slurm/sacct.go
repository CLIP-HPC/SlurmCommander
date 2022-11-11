package slurm

import (
	"github.com/pja237/slurmcommander-dev/internal/openapidb"
)

// SacctJobHist struct holds job history.
// Comes from unmarshalling sacct -A -S now-Xdays --json call.
type SacctJSON struct {
	Jobs []openapidb.Dbv0037Job
}

// This is to distinguish in Update() the return from the jobhisttab and jobdetails tab
type SacctSingleJobHist struct {
	Jobs []openapidb.Dbv0037Job
}
