package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type FileListKeymap struct {
	Add       key.Binding
	Remove    key.Binding
	Merge     key.Binding
	Save      key.Binding
	ShiftUp   key.Binding
	ShiftDown key.Binding
}

func MergeFullHelp() [][]key.Binding {
	return [][]key.Binding{
		{FileListKeys.Add, FileListKeys.Remove, FileListKeys.Merge},       // first column
		{FileListKeys.ShiftUp, FileListKeys.ShiftDown, FileListKeys.Save}, // second column
	}
}

var FileListKeys = FileListKeymap{
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
}
