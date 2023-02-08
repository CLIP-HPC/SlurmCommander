package abouttab

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/CLIP-HPC/SlurmCommander/internal/keybindings"
)

type Keys map[*key.Binding]bool

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.Up:       false,
	&keybindings.DefaultKeyMap.Down:     false,
	&keybindings.DefaultKeyMap.PageUp:   false,
	&keybindings.DefaultKeyMap.PageDown: false,
	&keybindings.DefaultKeyMap.Slash:    false,
	&keybindings.DefaultKeyMap.Info:     false,
	&keybindings.DefaultKeyMap.Enter:    false,
	&keybindings.DefaultKeyMap.Stats:    false,
	&keybindings.DefaultKeyMap.Count:    false,
	&keybindings.DefaultKeyMap.Params:   false,
	&keybindings.DefaultKeyMap.TimeRange:false,
}

func (k *Keys) SetupKeys() {
	for k, v := range KeyMap {
		k.SetEnabled(v)
	}
}
