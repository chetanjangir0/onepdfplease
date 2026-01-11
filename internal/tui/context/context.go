package context

import "github.com/chetanjangir0/onepdfplease/internal/tui/types"

type ProgramContext struct {
	MainContentHeight int
	TermWidth         int
	TermHeight        int
	Status            string
	StatusType        StatusType
	CurrentPage       types.Page
	// Config            *config.Config
}

type StatusType int

const (
	Error StatusType = iota
	Success
	None
)
