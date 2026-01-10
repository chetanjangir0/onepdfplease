package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type ListFileKeymap struct {
	Add       key.Binding
	Remove    key.Binding
	Merge     key.Binding
	Save      key.Binding
	ShiftUp   key.Binding
	ShiftDown key.Binding
	Help      key.Binding
}

func (k ListFileKeymap) ShortHelp() []key.Binding {
	return []key.Binding{k.Add, k.Remove, k.Merge, k.Save, k.Help}
}

func (k ListFileKeymap) FullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Add, k.Remove, k.Merge, k.Save}, // first column
		{k.ShiftUp, k.ShiftDown, k.Help},   // second column
	}
}

var ListFilesKeys = ListFileKeymap{
	Add: key.NewBinding(
		key.WithKeys("a"),
		key.WithHelp("a", "Add files"),
	),
	Remove: key.NewBinding(
		key.WithKeys("d"),
		key.WithHelp("d", "Remove files"),
	),
	Merge: key.NewBinding(
		key.WithKeys("m"),
		key.WithHelp("m", "Merge PDFs"),
	),
	Save: key.NewBinding(
		key.WithKeys("s"),
		key.WithHelp("s", "Save PDF"),
	),
	ShiftDown: key.NewBinding(
		key.WithKeys("J", "ctrl+down"),
		key.WithHelp("J/ctrl+down", "Shift Down"),
	),
	ShiftUp: key.NewBinding(
		key.WithKeys("K", "ctrl+up"),
		key.WithHelp("K/ctrl+up", "Shift Up"),
	),
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
}
