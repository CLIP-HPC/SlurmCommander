package main

import (
	"flag"
	"log"

	"github.com/pja237/slurmcommander-dev/internal/defaults"
)

type cmdflags struct {
	port    *uint
	tick    *uint
	version *bool
}

func (c *cmdflags) GetCmdFlags() {
	c.port = flag.Uint("p", defaults.SccPort, "port where to listen for SC requests")
	c.tick = flag.Uint("t", defaults.SccRefreshT, "time period, how often to fetch data from slurm")
	c.version = flag.Bool("v", false, "show version")
	flag.Parse()

}

func (c *cmdflags) DumpFlags() {
	log.Println("--------------------------------------------------------------------------------")
	log.Printf("Port: %d\n", *c.port)
	log.Printf("Version: %t\n", *c.version)
	log.Printf("Tick: %d\n", *c.tick)
	log.Println("--------------------------------------------------------------------------------")
}

func fetcherSacct() {

}

func main() {

	var cmd cmdflags

	log.Printf("- SCCache Start")

	// Parse CMDline parameters
	cmd.GetCmdFlags()
	cmd.DumpFlags()

	// Read config
	//  - skip, flags only

	// Setup logger
	//	- skip

	// Spin up scraper goroutine

	// Spin up bank

	// Spin up server listeners

	log.Printf("- SCCache End")
}
