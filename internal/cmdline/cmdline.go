package cmdline

import (
	"errors"
	"flag"
)

// CmdArgs holds currently supported command line parameters.
type CmdArgs struct {
	Version     *bool
	HistDays    *uint
	HistTimeout *uint
}

// NewCmdArgs return the CmdArgs structure built from command line parameters.
func NewCmdArgs() (*CmdArgs, error) {
	c := new(CmdArgs)

	c.Version = flag.Bool("v", false, "Display version")
	c.HistDays = flag.Uint("d", 7, "Jobs history fetch last N days")
	c.HistTimeout = flag.Uint("t", 30, "Job history fetch timeout, seconds")
	flag.Parse()
	if !flag.Parsed() {
		return nil, errors.New("failed to parse command line flags")
	}

	return c, nil
}
