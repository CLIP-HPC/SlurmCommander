package model

import (
	"log"

	"github.com/charmbracelet/bubbles/help"
	"github.com/CLIP-HPC/SlurmCommander/internal/config"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/abouttab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/clustertab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobdetailstab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobfromtemplate"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobhisttab"
	"github.com/CLIP-HPC/SlurmCommander/internal/model/tabs/jobtab"
)

const (
	tabJobs = iota
	tabJobHist
	tabJobDetails
	tabJobFromTemplate
	tabCluster
	tabAbout
)

// TODO: put this in model?
var tabs = []string{
	"Job Queue",
	"Job History",       // TODO: get this from sacct, even without json, then on Enter, goto: Job Details tab and fetch JSON there for a specific job?
	"Job Details",       // TODO: either show jobid textinput, or open this tab from Job History on selection
	"Job from Template", // TODO: devise sbatch templates and menus in this tab to fill them out
	"Cluster",
	"About",
}

type ActiveTabKeys interface {
	SetupKeys()
}

var tabKeys = []ActiveTabKeys{
	&jobtab.KeyMap,
	&jobhisttab.KeyMap,
	&jobdetailstab.KeyMap,
	&jobfromtemplate.KeyMap,
	&clustertab.KeyMap,
	&abouttab.KeyMap,
}

// TODO: in structures below:
// - make embedding and accessing leafs uniform (shorthand notation vs Full path)
type Model struct {
	Globals
	jobtab.JobTab
	jobhisttab.JobHistTab
	jobdetailstab.JobDetailsTab
	jobfromtemplate.JobFromTemplateTab
	clustertab.ClusterTab
}

type Globals struct {
	ActiveTab uint
	UpdateCnt uint64
	Debug     bool
	DebugMsg  string
	lastKey   string
	winW      int
	winH      int
	Log       *log.Logger
	Help      help.Model
	UserName  string
	UAccounts []string
	config.ConfigContainer
	ErrorMsg  error
	ErrorHelp string
	SizeErr   string
}
