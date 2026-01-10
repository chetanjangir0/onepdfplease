package context

type ProgramContext struct {
	MainContentHeight int
	TermWidth         int
	TermHeight        int
	Status            string
	StatusType        StatusType
	// Config            *config.Config
}

type StatusType int

const (
	Error StatusType  = iota
	Success
	None
)
