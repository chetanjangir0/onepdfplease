package main

// todo
// add progress bar

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal"
)

func main() {
	program := tea.NewProgram(internal.InitialModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		log.Fatal(err)
	}
}
