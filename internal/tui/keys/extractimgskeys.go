package keys

import "github.com/charmbracelet/bubbles/key"

func ExtractImgsFullHelp() [][]key.Binding {
	return [][]key.Binding{
		{FileListKeys.Add, FileListKeys.Remove},       // first column
		{FileListKeys.ShiftUp, FileListKeys.ShiftDown}, // second column
	}
}
