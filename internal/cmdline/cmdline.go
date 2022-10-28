package cmdline

import (
	"errors"
	"flag"
)

// CmdArgs holds currently supported command line parameters.
type CmdArgs struct {
	Version *bool
}

// NewCmdArgs return the CmdArgs structure built from command line parameters.
func NewCmdArgs() (*CmdArgs, error) {
	c := new(CmdArgs)

	c.Version = flag.Bool("v", false, "Display version")
	flag.Parse()
	if !flag.Parsed() {
		return nil, errors.New("failed to parse command line flags")
	}

	return c, nil
}
