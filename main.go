package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/cameronjpr/gaffer/internal/game"
	"github.com/cameronjpr/gaffer/internal/tui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	match := game.NewMatch()
	model := tui.NewModel(&match)

	p := tea.NewProgram(model, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}

	// writeMatchLog(finalModel.match)
	// After TUI exits, write match log
}

func runCmd(name string, arg ...string) {
	cmd := exec.Command(name, arg...)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func ClearTerminal() {
	switch runtime.GOOS {
	case "darwin":
		runCmd("clear")
	case "linux":
		runCmd("clear")
	case "windows":
		runCmd("cmd", "/c", "cls")
	default:
		runCmd("clear")
	}
}

/*
 * Pitch (zones)
 * Phases of play
 *
 *
 *
 * |---|---|---|
 * |   |   |   |
 * |---|---|---|
 * |   |   |   |
 * |---|---|---|
 * |   |   |   |
 * |---|---|---|
 * |   |   |   |
 * |---|---|---|
 */
