package listfiles 

import (
	"github.com/charmbracelet/bubbles/key"
)

type keyMap struct {
	add       key.Binding
	remove    key.Binding
	merge     key.Binding
	save      key.Binding
	shiftUp   key.Binding
	shiftDown key.Binding
	help      key.Binding
}

func (k keyMap) ShortHelp() []key.Binding {
	return []key.Binding{k.add, k.remove, k.merge, k.save, k.help}
}

func (k keyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.add, k.remove, k.merge, k.save}, // first column
		{k.shiftUp, k.shiftDown, k.help},   // second column
	}
}

var keys = keyMap{
	add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add files"),
	),
	remove: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Remove files"),
	),
	merge: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "Merge PDFs"),
	),
	save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Save PDF"),
	),
	shiftDown: key.NewBinding(
		key.WithKeys("J", "ctrl+down"),
		key.WithHelp("J/ctrl+down", "Shift Down"),
	),
	shiftUp: key.NewBinding(
		key.WithKeys("K", "ctrl+up"),
		key.WithHelp("K/ctrl+up", "Shift Up"),
	),
	help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}
