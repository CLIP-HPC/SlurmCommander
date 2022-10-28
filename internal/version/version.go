package version

import "fmt"

var (
	BuildVersion string
	BuildCommit  string
)

func DumpVersion() {
	fmt.Printf("----------------------------------------\n")
	fmt.Printf("Version: %s\n", BuildVersion)
	fmt.Printf("Build commit hash: %s\n", BuildCommit)
	fmt.Printf("----------------------------------------\n")
}
