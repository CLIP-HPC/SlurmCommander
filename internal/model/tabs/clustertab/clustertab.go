package clustertab

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/progress"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/pja237/slurmcommander/internal/keybindings"
	"github.com/pja237/slurmcommander/internal/slurm"
)

type JobClusterTab struct {
	SinfoTable    table.Model
	CpuBar        progress.Model
	MemBar        progress.Model
	Sinfo         slurm.SinfoJSON
	SinfoFiltered slurm.SinfoJSON
	Filter        textinput.Model
}

type Keys map[*key.Binding]bool

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.Up:       true,
	&keybindings.DefaultKeyMap.Down:     true,
	&keybindings.DefaultKeyMap.PageUp:   true,
	&keybindings.DefaultKeyMap.PageDown: true,
	&keybindings.DefaultKeyMap.Slash:    true,
	&keybindings.DefaultKeyMap.Info:     false,
	&keybindings.DefaultKeyMap.Enter:    false,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
}
