package jobfromtemplate

import (
	"github.com/CLIP-HPC/SlurmCommander/internal/keybindings"
	"github.com/charmbracelet/bubbles/key"
)

type Keys []*key.Binding

var KeyMap = Keys{
	&keybindings.DefaultKeyMap.Up,
	&keybindings.DefaultKeyMap.Down,
	&keybindings.DefaultKeyMap.Enter,
}

var EditorKeyMap = Keys{
	&keybindings.DefaultKeyMap.SaveSubmitJob,
	&keybindings.DefaultKeyMap.Escape,
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
