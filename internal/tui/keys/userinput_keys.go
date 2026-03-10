package keys

import (
	"github.com/charmbracelet/bubbles/key"
)

type UserInputKeymap struct {
	NxtField   key.Binding
	PrevField   key.Binding
	BtnPress    key.Binding
	ToggleField key.Binding
}

var UserInputKeys = UserInputKeymap{
	NxtField: key.NewBinding(
		key.WithKeys("down", "ctrl+n"),
		key.WithHelp("↓/ctrl+n", "Next field"),
	),
	PrevField: key.NewBinding(
		key.WithKeys("up", "ctrl+p"),
		key.WithHelp("↑/ctrl+p", "Prev field"),
	),
	BtnPress: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "Press button"),
	),
	ToggleField: key.NewBinding(
		key.WithKeys("left", "right", "h", "l"),
		key.WithHelp("left/right/h/l", "Toggle field value"),
	),
}

func UserInputFullHelp() [][]key.Binding {
	return [][]key.Binding{
		{UserInputKeys.NxtField, UserInputKeys.PrevField, UserInputKeys.ToggleField},
	}
}
