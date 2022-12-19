package jobfromtemplate

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/CLIP-HPC/SlurmCommander/internal/keybindings"
)

type Keys map[*key.Binding]bool

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.TtabSel:       true,
	&keybindings.DefaultKeyMap.Up:            true,
	&keybindings.DefaultKeyMap.Down:          true,
	&keybindings.DefaultKeyMap.PageUp:        false,
	&keybindings.DefaultKeyMap.PageDown:      false,
	&keybindings.DefaultKeyMap.Tab:           true,
	&keybindings.DefaultKeyMap.Slash:         false,
	&keybindings.DefaultKeyMap.Info:          false,
	&keybindings.DefaultKeyMap.Enter:         true,
	&keybindings.DefaultKeyMap.Quit:          true,
	&keybindings.DefaultKeyMap.SaveSubmitJob: false,
	&keybindings.DefaultKeyMap.Escape:        false,
	&keybindings.DefaultKeyMap.Stats:         false,
	&keybindings.DefaultKeyMap.Count:         false,
}

var EditorKeyMap = Keys{
	&keybindings.DefaultKeyMap.TtabSel:       false,
	&keybindings.DefaultKeyMap.Up:            false,
	&keybindings.DefaultKeyMap.Down:          false,
	&keybindings.DefaultKeyMap.PageUp:        false,
	&keybindings.DefaultKeyMap.PageDown:      false,
	&keybindings.DefaultKeyMap.Tab:           false,
	&keybindings.DefaultKeyMap.Slash:         false,
	&keybindings.DefaultKeyMap.Info:          false,
	&keybindings.DefaultKeyMap.Enter:         false,
	&keybindings.DefaultKeyMap.Quit:          false,
	&keybindings.DefaultKeyMap.SaveSubmitJob: true,
	&keybindings.DefaultKeyMap.Escape:        true,
	&keybindings.DefaultKeyMap.Stats:         false,
	&keybindings.DefaultKeyMap.Count:         false,
}

func (k *Keys) SetupKeys() {
	//for k, v := range KeyMap {
	for k, v := range *k {
		k.SetEnabled(v)
	}
}
