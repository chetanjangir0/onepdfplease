package keys

import (
	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

type KeyMap struct {
	Page types.Page
	Help key.Binding
	Quit key.Binding
}

var keys = &KeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q", "esc"),
		key.WithHelp("q", "quit"),
	),
}

func CreateKeyMapForPage(page types.Page) help.KeyMap {
	keys.Page = page
	return keys
}

func (k KeyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.Help}
}

func (k KeyMap) FullHelp() [][]key.Binding {
	var additionalKeys [][]key.Binding
	switch k.Page {
	case types.MergePage:
		additionalKeys = MergeFullHelp()
	case types.MenuPage:
		additionalKeys = MenuFullHelp()
	default:
		additionalKeys = MenuFullHelp()
	}

	allKeys := append(additionalKeys, k.GlobalKeys())
	return allKeys
}

func (k KeyMap) GlobalKeys() []key.Binding {
	return []key.Binding{
		k.Help,
		k.Quit,
	}
}
