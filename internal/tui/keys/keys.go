package keys

import (
	"slices"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/key"
	"github.com/chetanjangir0/onepdfplease/internal/tui/types"
)

type KeyMap struct {
	Page    types.Page
	Help    key.Binding
	Back    key.Binding
	Quit    key.Binding
	NxtTab  key.Binding
	PrevTab key.Binding
}

var Keys = &KeyMap{
	Help: key.NewBinding(
		key.WithKeys("?"),
		key.WithHelp("?", "toggle help"),
	),
	Back: key.NewBinding(
		key.WithKeys("esc"),
		key.WithHelp("Esc", "Go Back"),
	),
	Quit: key.NewBinding(
		key.WithKeys("q"),
		key.WithHelp("q", "quit"),
	),
	NxtTab: key.NewBinding(
		key.WithKeys("tab"),
		key.WithHelp("Tab", "Next tab"),
	),
	PrevTab: key.NewBinding(
		key.WithKeys("shift+tab"),
		key.WithHelp("shift+tab", "Previous tab"),
	),
}

func CreateKeyMapForPage(page types.Page) help.KeyMap {
	Keys.Page = page
	return Keys
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
	case types.EncryptPage:
		additionalKeys = EncryptFullHelp()
	case types.SplitPage:
		additionalKeys = SplitFullHelp()
	case types.DecryptPage:
		additionalKeys = DecryptFullHelp()
	case types.Img2PdfPage:
		additionalKeys = Img2pdfFullHelp()
	case types.ExtractImgsPage:
		additionalKeys = ExtractImgsFullHelp()
	default:
		additionalKeys = k.GlobalFullHelp()
	}

	allKeys := slices.Concat(additionalKeys[:], k.GlobalKeys()[:])
	return allKeys
}

func (k KeyMap) GlobalKeys() [][]key.Binding {
	return [][]key.Binding{
		{k.NxtTab, k.PrevTab, k.Help, k.Back},
		{k.Quit},
	}
}

func (k KeyMap) GlobalFullHelp() [][]key.Binding {
	return [][]key.Binding{
		{k.Help},
		{k.Back},
		{k.Quit},
		{k.NxtTab},
		{k.PrevTab},
	}
}
