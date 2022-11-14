/*
Package config implements the ConfigContainer structure and accompanying methods.
It holds the configuration data for all utilities.
Configuration file format is the same for all.
*/
package config

import (
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/BurntSushi/toml"
)

type ConfigContainer struct {
	Prefix   string            // if this is set, then we prepend this path to all commands
	Binpaths map[string]string `json:"binpaths"` // else, we specify one by one
	Tick     uint
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
		errNo    uint
	)
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	cfgPaths = []string{"/etc/scom/scom.conf", home + "/scom/scom.conf"}
	for _, v := range cfgPaths {
		log.Printf("Trying conf file: %s\n", v)
		f, err := os.ReadFile(v)
		if err != nil {
			//return err
			errNo++
			continue
		}

		err = toml.Unmarshal(f, cc)

		if err != nil {
			cc.testNsetBinPaths()
			return err
		}
	}

	if cc.Tick == 0 {
		// set default Tick
		cc.Tick = 3
	}
	cc.testNsetBinPaths()

	if errNo == 2 {
		return errors.New("/etc/scom/scom.conf or $HOME/scom/scom.conf NOT FOUND")
	}

	return nil
}

func (cc *ConfigContainer) testNsetBinPaths() error {

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

	return nil
}

func (cc *ConfigContainer) DumpConfig() string {

	return fmt.Sprintf("Configuration: %#v\n", cc)

}
