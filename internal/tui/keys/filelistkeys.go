package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type FileListKeymap struct {
	Add            key.Binding
	Remove         key.Binding
	ShiftUp        key.Binding
	ShiftDown      key.Binding
	QuitFilepicker key.Binding
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
	ShiftDown: key.NewBinding(
		key.WithKeys("J", "ctrl+down"),
		key.WithHelp("J/ctrl+down", "Shift Down"),
	),
	ShiftUp: key.NewBinding(
		key.WithKeys("K", "ctrl+up"),
		key.WithHelp("K/ctrl+up", "Shift Up"),
	),
	QuitFilepicker: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("esc", "Go Back"),
	),
}

func FileListFullHelp() [][]key.Binding {
	return [][]key.Binding{
		{FileListKeys.Add, FileListKeys.Remove, FileListKeys.ShiftUp, FileListKeys.ShiftDown},
	}

}
