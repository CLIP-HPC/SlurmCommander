package defaults

const (
	TickMin = 3 // minimal time in seconds between that can be set in config file. If not set or less then, Set to this value.
	HistTimeout = 30 // in seconds; must be >= 1
	HistStart = "now-7days" // must be >= 1

	AppName = "scom"

	EnvConfVarName = "SCOM_CONF"
	ConfFileName   = "scom.conf"
	SiteConfDir    = "/etc/" + AppName + "/"
	SiteConfFile   = SiteConfDir + ConfFileName

	TemplatesDir        = SiteConfDir + "templates"
	TemplatesSuffix     = ".sbatch"
	TemplatesDescSuffix = ".desc"
)

var (
	// default paths
	BinPaths = map[string]string{
		"sacct":    "/bin/sacct",
		"sstat":    "/bin/sstat",
		"sinfo":    "/bin/sinfo",
		"squeue":   "/bin/squeue",
		"sbatch":   "/bin/sbatch",
		"scancel":  "/bin/scancel",
		"scontrol": "/bin/scontrol",
		"sacctmgr": "/bin/sacctmgr",
	}
)
