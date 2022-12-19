package jobhisttab

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/CLIP-HPC/SlurmCommander/internal/command"
	"github.com/CLIP-HPC/SlurmCommander/internal/config"
)

var (
	cc                   config.ConfigContainer
	SacctHistCmdSwitches = []string{"-n", "--json"}
)

func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

type JobHistTabMsg struct {
	HistFetchFail bool
	SacctJSON
}

func GetSacctHist(uaccs string, d uint, t uint, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var (
			jht JobHistTabMsg
		)

		// TODO: Setup command line switch to pick how many days of sacct to fetch in case of massive runs.

		l.Printf("GetSacctHist(%q) start: days %d, timeout: %d\n", uaccs, d, t)

		// setup context with 5 second timeout
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
		defer cancel()

		// prepare command
		cmd := cc.Binpaths["sacct"]
		start := fmt.Sprintf("now-%ddays", d)
		switches := append(SacctHistCmdSwitches, "-S", start)
		switches = append(switches, "-A", uaccs)

		l.Printf("EXEC: %q %q\n", cmd, switches)
		out, err := exec.CommandContext(ctx, cmd, switches...).CombinedOutput()
		if err != nil {
			l.Printf("Error exec sacct: %q\n", err)
			// set error, return.
			// TODO: see how to fit this with the new commands-return-command.ErrorMsg pattern
			jht.HistFetchFail = true
			return jht
		}
		l.Printf("EXEC returned: %d bytes\n", len(out))

		err = json.Unmarshal(out, &jht.SacctJSON)
		if err != nil {
			//jht.HistFetchFail = true
			l.Printf("Error unmarshall: %q\n", err)
			return command.ErrorMsg{
				From:    "GetSacctHist",
				ErrHelp: "sacct JSON failed to parse, note your slurm version and open an issue with us here: https://github.com/CLIP-HPC/SlurmCommander/issues/new/choose",
				OrigErr: err,
			}
			//return jht
		}

		jht.HistFetchFail = false
		l.Printf("Unmarshalled %d jobs from hist\n", len(jht.SacctJSON.Jobs))

		return jht
	}
}
