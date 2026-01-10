package context

type ProgramContext struct {
	MainContentHeight int
	TermWidth         int
	TermHeight        int
	// Config            *config.Config
	Error error
}
