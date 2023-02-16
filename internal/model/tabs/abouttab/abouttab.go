package abouttab

import (
	"github.com/charmbracelet/bubbles/key"
)

type Keys []*key.Binding

var KeyMap = Keys{}

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
