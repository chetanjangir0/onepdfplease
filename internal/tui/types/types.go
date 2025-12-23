package types

type Page int

const (
	MenuPage Page = iota
	MergePage
	SplitPage
)

type NavigateMsg struct {
	Page Page
}
