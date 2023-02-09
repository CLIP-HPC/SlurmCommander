package slurm

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/CLIP-HPC/SlurmCommander/internal/openapi"
)

type SinfoJSON struct {
	Nodes []openapi.V0039Node
}

var gpuGresPattern = regexp.MustCompile(`gpu:(.*:)?(\d+)\(.*`)

func ParseGRES(line string) *int {
	value := 0

	gres := strings.Split(line, ",")
	for _, g := range gres {
		if !strings.HasPrefix(g, "gpu:") {
			continue
		}

		matches := gpuGresPattern.FindStringSubmatch(g)
		if len(matches) == 3 {
			v, _ := strconv.Atoi(matches[2])
			value += v
		}
	}

	return &value
}

type GresMap map[string]int

func ParseGRESAll(line string) *GresMap {
	var gmap GresMap = make(GresMap)

	gres := strings.Split(line, ",")
	for _, g := range gres {
		if !strings.HasPrefix(g, "gpu:") {
			continue
		}

		matches := gpuGresPattern.FindStringSubmatch(g)
		if len(matches) == 3 {
			v, _ := strconv.Atoi(matches[2])
			gmap[matches[1]] += v
		}
	}

	return &gmap
}
