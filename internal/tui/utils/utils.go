package utils 

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// func (m *model) SetCurrentPage(p Page) {
// 	m.currentPage = p
// }

func Log(message string) {
	f, err := tea.LogToFile("./debug.log", message)
	if err != nil {
		fmt.Println("Fatal:", err)
		os.Exit(1)
	}
	defer f.Close()
}
