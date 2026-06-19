package main

import (
	"fmt"
	"log"

	"github.com/WeepingDogel/arch-server-installation-tui/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	model := tui.New()
	program := tea.NewProgram(model, tea.WithAltScreen())

	if _, err := program.Run(); err != nil {
		log.Fatalf("Error running installer: %v", err)
	}

	fmt.Println("Arch Linux Server installation complete. Reboot to start your new server!")
}