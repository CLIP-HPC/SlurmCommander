package keybindings

import (
	"github.com/charmbracelet/bubbles/key"
)

type KeyMap struct {
	TtabSel       key.Binding
	Up            key.Binding
	Down          key.Binding
	Quit          key.Binding
	PageUp        key.Binding
	PageDown      key.Binding
	Tab           key.Binding
	Slash         key.Binding
	Info          key.Binding
	Enter         key.Binding
	SaveSubmitJob key.Binding
	Escape        key.Binding
	Stats         key.Binding
}

// TODO: add shift+tab
var DefaultKeyMap = KeyMap{
	// TODO: combine tab selection keys into one and distinguish by Key.Value?
	TtabSel: key.NewBinding(
		key.WithKeys("1", "2", "3", "4", "5", "6"),
		key.WithHelp("1-6", "GoTo Tab"),
	),
	Stats: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Show Statistics"),
	),
	Up: key.NewBinding(
		key.WithKeys("k", "up"),        // actual keybindings
		key.WithHelp("↑/k", "Move up"), // corresponding help text
	),
	Down: key.NewBinding(
		key.WithKeys("j", "down"),
		key.WithHelp("↓/j", "Move down"),
	),
	PageUp: key.NewBinding(
		key.WithKeys("b", "pgup"),
		key.WithHelp("b/pgup", "Page Up"),
	),
	PageDown: key.NewBinding(
		key.WithKeys("f", "pgdown"),
		key.WithHelp("f/pgdn", "Page Down"),
	),
	Tab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("tab", "Cycle tabs"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "ctrl+c"),
		key.WithHelp("q", "Quit scom"),
	),
	Slash: key.NewBinding(
		key.WithKeys("/"),
		key.WithHelp("/", "Filter table"),
	),
	Info: key.NewBinding(
		key.WithKeys("i"),
		key.WithHelp("i", "Info"),
	),
	Enter: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Select entry"),
	),
	SaveSubmitJob: key.NewBinding(
		key.WithKeys("ctrl+s"),
		key.WithHelp("ctrl+s", "Save and Submit the job script"),
		key.WithDisabled(),
	),
	Escape: key.NewBinding(
		key.WithKeys("Esc"),
		key.WithHelp("Esc", "Exit without saving"),
		key.WithDisabled(),
	),
}

func (km KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{
		km.TtabSel,
		km.Up,
		km.Down,
		km.PageUp,
		km.PageDown,
		km.Tab,
		km.Slash,
		km.Info,
		km.Stats,
		km.Enter,
		km.Quit,
		km.SaveSubmitJob,
		km.Escape,
	}
}

func (km KeyMap) FullHelp() [][]key.Binding {
	// TODO: this...
	// MoreHelp returns an extended group of help items, grouped by columns.
	// The help bubble will render the help in the order in which the help
	// items are returned here.
	return nil
}
