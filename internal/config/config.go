/*
Package config implements the ConfigContainer structure and accompanying methods.
It holds the configuration data for all utilities.
Configuration file format is the same for all.
*/
package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type ConfigContainer struct {
	Prefix       string            // if this is set, then we prepend this path to all commands
	Binpaths     map[string]string // else, we specify one by one
	Tick         uint
	TemplateDirs []string
}

func NewConfigContainer() *ConfigContainer {
	return new(ConfigContainer)
}

func (cc *ConfigContainer) GetTick() time.Duration {
	return time.Duration(cc.Tick)
}

// Read & unmarshall configuration from 'name' file into configContainer structure
func (cc *ConfigContainer) GetConfig() error {
	var (
		cfgPaths []string
	)

	home, err := os.UserHomeDir()
	if err != nil {
		log.Printf("Conf: FAILED getting users $HOME %s\n", err)
		cfgPaths = []string{"/etc/scom/scom.conf"}
	} else {
		cfgPaths = []string{"/etc/scom/scom.conf", home + "/scom/scom.conf"}
	}

	// TODO: JFT this NEEDS rework, not good.
	for _, v := range cfgPaths {
		log.Printf("Trying conf file: %s\n", v)
		f, err := os.ReadFile(v)
		if err != nil {
			log.Printf("Conf: FAILED reading %s\n", v)
			continue
		}

		err = toml.Unmarshal(f, cc)
		if err != nil {
			log.Printf("Conf: FAILED unmarshalling %s with %s\n", v, err)
		}
	}

	// Here we test config limits and set them.
	// Also fill out unset config params.

	// if unset (==0) or less then 3, set to default
	if cc.Tick < 3 {
		// set default Tick
		cc.Tick = 3
	}
	cc.testNsetBinPaths()
	cc.testNsetTemplateDirs()

	// We don't return error since we set sane defaults and
	// errors arising from bad config should be handled in app.
	// for now leave signature as-is, later remove error return

	return nil
}

func (cc *ConfigContainer) testNsetTemplateDirs() {
	// default, always there:	"/etc/slurmcommander/templates",
	if cc.TemplateDirs == nil {
		// Nothing set from config files
		cc.TemplateDirs = append(cc.TemplateDirs, "/etc/slurmcommander/templates")
	} else {
		// Something exists from config, can be site-wide OR user-conf
		// QUESTION: should we do anything about it? prepend /etc/... one? or leave it as-is?
		// For now, we don't touch it.
	}

}

func (cc *ConfigContainer) testNsetBinPaths() {

	if cc.Binpaths == nil {
		cc.Binpaths = make(map[string]string)
	}

	// default paths
	defaultpaths := map[string]string{
		"sacct":    "/bin/sacct",
		"sstat":    "/bin/sstat",
		"sinfo":    "/bin/sinfo",
		"squeue":   "/bin/squeue",
		"sbatch":   "/bin/sbatch",
		"scancel":  "/bin/scancel",
		"scontrol": "/bin/scontrol",
		"sacctmgr": "/bin/sacctmgr",
	}

	for key, path := range defaultpaths {
		if val, exists := cc.Binpaths[key]; !exists || val == "" {
			if cc.Prefix != "" {
				// prefix is set, prepend it
				cc.Binpaths[key] = cc.Prefix + "/" + key
			} else {
				cc.Binpaths[key] = path
			}
		}
	}

}

func (cc *ConfigContainer) DumpConfig() string {

	return fmt.Sprintf("Configuration: %#v\n", cc)

}
