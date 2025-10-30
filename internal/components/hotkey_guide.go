package components

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// HotkeyBinding represents a single keyboard shortcut
type HotkeyBinding struct {
	Key         string
	Description string
}

// HotkeyGuide renders a footer with contextual hotkey bindings
func HotkeyGuide(width int, bindings []HotkeyBinding) string {
	if len(bindings) == 0 {
		return ""
	}

	// Style for keys
	keyStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170"))

	// Style for descriptions
	descStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("252"))

	// Build hotkey items
	items := make([]string, len(bindings))
	for i, binding := range bindings {
		items[i] = keyStyle.Render(binding.Key) + descStyle.Render(": "+binding.Description)
	}

	// Join with separator
	separator := lipgloss.NewStyle().
		Foreground(lipgloss.Color("240")).
		Render(" • ")

	hotkeyText := strings.Join(items, separator)

	// Center the hotkeys
	centered := lipgloss.NewStyle().
		Width(width).
		Align(lipgloss.Center).
		Render(hotkeyText)

	// Wrap in a bordered style
	footer := lipgloss.NewStyle().
		Width(width).
		Border(lipgloss.NormalBorder(), true, false, false, false).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 0).
		Render(centered)

	return footer
}
