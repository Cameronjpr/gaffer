package components

import (
	"github.com/charmbracelet/lipgloss"
)

func ThreeColumnLayout(width int, c1, c2, c3 string) string {
	colWidth := width / 3

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render(c1),
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render(c2),
		lipgloss.NewStyle().Padding(2).Width(colWidth).Render(c3),
	)
}

// ModalConfig contains configuration for modal rendering
type ModalConfig struct {
	Width        int
	Height       int
	MarginX      int // Horizontal margin from screen edges
	MarginY      int // Vertical margin from screen edges
	BorderStyle  lipgloss.Border
	BorderColor  string
	ShowBorder   bool
	Background   string
	Foreground   string
	Title        string
	ContentAlign lipgloss.Position
}

// DefaultModalConfig returns sensible defaults for modal rendering
func DefaultModalConfig(width, height int) ModalConfig {
	return ModalConfig{
		Width:        width,
		Height:       height,
		MarginX:      8,
		MarginY:      4,
		BorderStyle:  lipgloss.RoundedBorder(),
		BorderColor:  "240",
		ShowBorder:   true,
		Background:   "",
		Foreground:   "",
		Title:        "",
		ContentAlign: lipgloss.Center,
	}
}

// Modal renders content in a centered modal with chunky margins
// This creates an "overlay" effect that sits in from the screen edges
func Modal(config ModalConfig, content string) string {
	// Calculate available space after margins
	availableWidth := config.Width - (config.MarginX * 2)
	availableHeight := config.Height - (config.MarginY * 2)

	// Ensure we don't go negative
	if availableWidth < 20 {
		availableWidth = 20
	}
	if availableHeight < 5 {
		availableHeight = 5
	}

	// Build the content style
	contentStyle := lipgloss.NewStyle().
		Width(availableWidth).
		Height(availableHeight).
		Align(config.ContentAlign, lipgloss.Center)

	// Add border if requested
	if config.ShowBorder {
		contentStyle = contentStyle.
			Border(config.BorderStyle).
			BorderForeground(lipgloss.Color(config.BorderColor))
	}

	// Add background/foreground if provided
	if config.Background != "" {
		contentStyle = contentStyle.Background(lipgloss.Color(config.Background))
	}
	if config.Foreground != "" {
		contentStyle = contentStyle.Foreground(lipgloss.Color(config.Foreground))
	}

	// Add title if provided
	var finalContent string
	if config.Title != "" {
		// Use custom title style or default
		titleStyle := lipgloss.NewStyle().
			Bold(true).
			Padding(0, 1).
			Align(lipgloss.Center)

		titleRendered := titleStyle.Render(config.Title)
		finalContent = lipgloss.JoinVertical(lipgloss.Center, titleRendered, "", content)
	} else {
		finalContent = content
	}

	// Render the modal content
	modalContent := contentStyle.Render(finalContent)

	// Center the modal on the full screen with margins
	return lipgloss.Place(
		config.Width,
		config.Height,
		lipgloss.Center,
		lipgloss.Center,
		modalContent,
	)
}

// SimpleModal is a convenience wrapper for Modal with default settings
func SimpleModal(width, height int, title, content string) string {
	config := DefaultModalConfig(width, height)
	config.Title = title
	return Modal(config, content)
}

// Panel renders content in a padded container with optional border
// Useful for creating consistent panels throughout the UI
type PanelConfig struct {
	Width       int
	Height      int
	PaddingX    int
	PaddingY    int
	Border      bool
	BorderStyle lipgloss.Border
	BorderColor string
	Background  string
	Foreground  string
	Align       lipgloss.Position
}

func DefaultPanelConfig(width, height int) PanelConfig {
	return PanelConfig{
		Width:       width,
		Height:      height,
		PaddingX:    2,
		PaddingY:    1,
		Border:      true,
		BorderStyle: lipgloss.NormalBorder(),
		BorderColor: "240",
		Align:       lipgloss.Left,
	}
}

func Panel(config PanelConfig, content string) string {
	style := lipgloss.NewStyle().
		Width(config.Width).
		Padding(config.PaddingY, config.PaddingX).
		Align(config.Align)

	if config.Border {
		style = style.
			Border(config.BorderStyle).
			BorderForeground(lipgloss.Color(config.BorderColor))
	}

	if config.Background != "" {
		style = style.Background(lipgloss.Color(config.Background))
	}
	if config.Foreground != "" {
		style = style.Foreground(lipgloss.Color(config.Foreground))
	}

	// Apply height if specified (after padding/border)
	if config.Height > 0 {
		style = style.Height(config.Height)
	}

	return style.Render(content)
}

// Grid creates a responsive grid layout
// Useful for consistent spacing of elements
func Grid(width int, columns int, items []string) string {
	if columns <= 0 {
		columns = 1
	}
	if len(items) == 0 {
		return ""
	}

	colWidth := width / columns
	var rows []string
	var currentRow []string

	for i, item := range items {
		currentRow = append(currentRow, lipgloss.NewStyle().Width(colWidth).Render(item))

		// End of row or last item
		if (i+1)%columns == 0 || i == len(items)-1 {
			rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, currentRow...))
			currentRow = []string{}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, rows...)
}

// Stack creates a vertical stack of items with consistent spacing
func Stack(items []string, spacing int) string {
	if len(items) == 0 {
		return ""
	}

	// Add spacing between items
	var spacedItems []string
	for i, item := range items {
		spacedItems = append(spacedItems, item)
		// Add spacing except after last item
		if i < len(items)-1 && spacing > 0 {
			for j := 0; j < spacing; j++ {
				spacedItems = append(spacedItems, "")
			}
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, spacedItems...)
}

// Centered centers content within a given width and height
func Centered(width, height int, content string) string {
	return lipgloss.Place(width, height, lipgloss.Center, lipgloss.Center, content)
}

// ScreenSection divides the screen into sections with consistent spacing.
// Useful for creating main layout areas (header, content, footer)
type ScreenSection struct {
	Height  int // Fixed height, or 0 for flexible
	Content string
}

func ScreenLayout(totalHeight int, sections []ScreenSection) string {
	// Calculate flexible height
	fixedHeight := 0
	flexCount := 0
	for _, section := range sections {
		if section.Height > 0 {
			fixedHeight += section.Height
		} else {
			flexCount++
		}
	}

	flexHeight := 0
	if flexCount > 0 {
		remaining := totalHeight - fixedHeight
		if remaining > 0 {
			flexHeight = remaining / flexCount
		}
	}

	// Render sections with proper heights
	var rendered []string
	for _, section := range sections {
		height := section.Height
		if height == 0 {
			height = flexHeight
		}
		// Ensure content takes up the full calculated height
		sectionContent := lipgloss.NewStyle().
			Height(height).
			Render(section.Content)
		rendered = append(rendered, sectionContent)
	}

	return lipgloss.JoinVertical(lipgloss.Top, rendered...)
}
