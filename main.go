package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cameronjpr/gaffer/internal/db"
	"github.com/cameronjpr/gaffer/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize database
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Println("Error getting home directory:", err)
		os.Exit(1)
	}

	dbPath := filepath.Join(homeDir, ".gaffer", "gaffer.db")
	database, err := db.InitDB(dbPath)
	if err != nil {
		fmt.Println("Error initializing database:", err)
		os.Exit(1)
	}
	defer database.Close()

	// Seed database with clubs from JSON
	if err := db.SeedDatabase(database, "clubs.json"); err != nil {
		fmt.Println("Error seeding database:", err)
		os.Exit(1)
	}

	// Create queries
	queries := db.New(database)

	model := tui.NewModel(queries)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
