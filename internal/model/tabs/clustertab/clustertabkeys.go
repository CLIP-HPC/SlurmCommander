package clustertab

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/CLIP-HPC/SlurmCommander/internal/keybindings"
)

type Keys []*key.Binding

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.Up,
	&keybindings.DefaultKeyMap.Down,
	&keybindings.DefaultKeyMap.PageUp,
	&keybindings.DefaultKeyMap.PageDown,
	&keybindings.DefaultKeyMap.Slash,
	&keybindings.DefaultKeyMap.Stats,
	&keybindings.DefaultKeyMap.Count,
}

func (ky *Keys) SetupKeys() {
	for _, k := range *ky {
		k.SetEnabled(true)
	}
}

func (ky *Keys) DisableKeys() {
	for _, k := range *ky {
		k.SetEnabled(false)
	}
}
