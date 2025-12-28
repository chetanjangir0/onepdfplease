package main

// todo
// add progress bar
// add stack like structure for easier file navigation

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/chetanjangir0/onepdfplease/internal/tui"
)

func main() {
	f, err := tea.LogToFile("./debug.log", "DEBUG:")
	if err != nil {
		fmt.Println("Fatal:", err)
		os.Exit(1)
	}
	defer f.Close()

	program := tea.NewProgram(tui.InitialModel(), tea.WithAltScreen())
	if _, err := program.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
