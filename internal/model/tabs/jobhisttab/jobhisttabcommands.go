package jobhisttab

import (
	"context"
	"strings"
	"encoding/json"
	"log"
	"os/exec"
	"time"

	"github.com/CLIP-HPC/SlurmCommander/internal/command"
	"github.com/CLIP-HPC/SlurmCommander/internal/config"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	cc config.ConfigContainer
	SacctHistCmdSwitches = []string{"-n", "--json"}
)

func NewCmdCC(config config.ConfigContainer) {
	cc = config
}

type JobHistTabMsg struct {
	HistFetchFail bool
	SacctJSON
}

func GetSacctHist(uaccs string, start string, end string, t uint, l *log.Logger) tea.Cmd {
	return func() tea.Msg {
		var (
			jht JobHistTabMsg
		)

		l.Printf("GetSacctHist(%q) start: %s, end: %s, timeout: %d\n", uaccs, start, end, t)

		// setup context with user set timeout in seconds
		ctx, cancel := context.WithTimeout(context.Background(), time.Duration(t)*time.Second)
		defer cancel()

		// prepare command
		cmd := cc.Binpaths["sacct"]
		switches := append(SacctHistCmdSwitches, "-A", uaccs)
		if strings.TrimSpace(start) != "" {
			switches = append(switches, "-S", start)
		}
		if strings.TrimSpace(end) != "" {
			switches = append(switches, "-E", end)
		}

		l.Printf("EXEC: %q %q\n", cmd, switches)
		out, err := exec.CommandContext(ctx, cmd, switches...).Output()
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
			l.Printf("Error unmarshall: %q\n", err)
			return command.ErrorMsg{
				From:    "GetSacctHist",
				ErrHelp: "sacct JSON failed to parse, note your slurm version and open an issue with us here: https://github.com/CLIP-HPC/SlurmCommander/issues/new/choose",
				OrigErr: err,
			}
		}

		jht.HistFetchFail = false
		l.Printf("Unmarshalled %d jobs from hist\n", len(jht.SacctJSON.Jobs))

		return jht
	}
}
