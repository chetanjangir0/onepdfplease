package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type MenuKeyMap struct {
	CursorUp    key.Binding
	CursorDown  key.Binding
	GoToStart   key.Binding
	GoToEnd     key.Binding
	Filter      key.Binding
	ClearFilter key.Binding
}

func MenuFullHelp() [][]key.Binding {
	return [][]key.Binding{
		{MenuKeys.CursorDown, MenuKeys.CursorUp, MenuKeys.GoToStart, MenuKeys.GoToEnd},
	}
}

var MenuKeys = MenuKeyMap{
	// Browsing.
	CursorUp: key.NewBinding(
		key.WithKeys("up", "k"),
		key.WithHelp("↑/k", "up"),
	),
	CursorDown: key.NewBinding(
		key.WithKeys("down", "j"),
		key.WithHelp("↓/j", "down"),
	),
	GoToStart: key.NewBinding(
		key.WithKeys("home", "g"),
		key.WithHelp("g/home", "go to start"),
	),
	GoToEnd: key.NewBinding(
		key.WithKeys("end", "G"),
		key.WithHelp("G/end", "go to end"),
	),
	// Filter: key.NewBinding(
	// 	key.WithKeys("/"),
	// 	key.WithHelp("/", "filter"),
	// ),
	// ClearFilter: key.NewBinding(
	// 	key.WithKeys("esc"),
	// 	key.WithHelp("esc", "clear filter"),
	// ),

}
