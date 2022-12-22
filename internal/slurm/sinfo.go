package slurm

import (
	"log"
	"regexp"
	"strconv"
	"strings"

	"github.com/CLIP-HPC/SlurmCommander/internal/openapi"
)

type SinfoJSON struct {
	Nodes []openapi.V0039Node
}

//var gpuGresPattern = regexp.MustCompile(`^gpu\:([^\:]+)\:?(\d+)?`)
var gpuGresPattern = regexp.MustCompile(`gpu:(.*:)?(\d+)\(.*`)

func ParseGRES(line string) *int {
	value := 0

	log.Printf("GOT: %q\n", line)
	gres := strings.Split(line, ",")
	for _, g := range gres {
		if !strings.HasPrefix(g, "gpu:") {
			continue
		}

		matches := gpuGresPattern.FindStringSubmatch(g)
		if len(matches) == 3 {
			value, _ = strconv.Atoi(matches[2])
		}
		//if len(matches) == 3 {
		//	if matches[2] != "" {
		//		value, _ = strconv.Atoi(matches[2])
		//	} else {
		//		value, _ = strconv.Atoi(matches[1])
		//	}
		//}
	}

	return &value
}
