package tui

import (
	"fmt"
	"io"
	"strings"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// SubstitutionModel handles the substitution UI with two side-by-side lists
type SubstitutionModel struct {
	width         int
	height        int
	fieldedList   list.Model
	benchList     list.Model
	focusedList   int // 0 = fielded, 1 = bench
	selectedXI    *domain.MatchPlayerParticipant
	selectedBench *domain.MatchPlayerParticipant
}

// NewSubstitutionModel creates a new substitution model
func NewSubstitutionModel(currentXI, bench []*domain.MatchPlayerParticipant) *SubstitutionModel {
	// Create fielded players list (Starting XI)
	fieldedItems := make([]list.Item, len(currentXI))
	for i, p := range currentXI {
		fieldedItems[i] = playerItem{player: p, isBench: false}
	}

	fieldedDelegate := playerDelegate{}
	fieldedList := list.New(fieldedItems, fieldedDelegate, 0, 0)
	fieldedList.SetShowTitle(false)
	fieldedList.SetShowStatusBar(false)
	fieldedList.SetShowHelp(false)
	fieldedList.SetFilteringEnabled(false)

	// Create bench players list
	benchItems := make([]list.Item, len(bench))
	for i, p := range bench {
		benchItems[i] = playerItem{player: p, isBench: true}
	}

	benchDelegate := playerDelegate{}
	benchList := list.New(benchItems, benchDelegate, 0, 0)
	benchList.SetShowTitle(false)
	benchList.SetShowStatusBar(false)
	benchList.SetShowHelp(false)
	benchList.SetFilteringEnabled(false)

	return &SubstitutionModel{
		fieldedList: fieldedList,
		benchList:   benchList,
		focusedList: 0, // Start with fielded list focused
	}
}

func (m *SubstitutionModel) Init() tea.Cmd {
	return nil
}

func (m *SubstitutionModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Each list takes roughly half the modal width
		listWidth := (msg.Width / 2) - 6
		// Reserve space for: title (1), empty line (1), status (2), instructions (2), borders/padding (~6)
		listHeight := msg.Height - 12
		m.fieldedList.SetSize(listWidth, listHeight)
		m.benchList.SetSize(listWidth, listHeight)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "tab", "left", "right":
			// Switch focus between lists
			m.focusedList = 1 - m.focusedList
			return m, nil

		case " ":
			// Spacebar toggles selection in the focused list
			if m.focusedList == 0 {
				// Fielded list
				selected := m.fieldedList.SelectedItem()
				if item, ok := selected.(playerItem); ok {
					// Toggle selection: if same player, deselect; otherwise select
					if m.selectedXI != nil && m.selectedXI.Player == item.player.Player {
						m.selectedXI = nil
					} else {
						m.selectedXI = item.player
					}
					// Always update delegate to reflect new state
					m.fieldedList.SetDelegate(playerDelegate{
						selectedPlayer: m.selectedXI,
					})
				}
			} else {
				// Bench list
				selected := m.benchList.SelectedItem()
				if item, ok := selected.(playerItem); ok {
					// Toggle selection: if same player, deselect; otherwise select
					if m.selectedBench != nil && m.selectedBench.Player == item.player.Player {
						m.selectedBench = nil
					} else {
						m.selectedBench = item.player
					}
					// Always update delegate to reflect new state
					m.benchList.SetDelegate(playerDelegate{
						selectedPlayer: m.selectedBench,
					})
				}
			}
			return m, nil

		case "enter":
			// Enter only works when both players are selected
			if m.selectedXI != nil && m.selectedBench != nil {
				return m, func() tea.Msg {
					return executeSubstitutionMsg{
						playerOut: m.selectedXI,
						playerIn:  m.selectedBench,
					}
				}
			}
			return m, nil

		case "esc":
			// Clear selections
			m.selectedXI = nil
			m.selectedBench = nil
			// Update delegates
			m.fieldedList.SetDelegate(playerDelegate{})
			m.benchList.SetDelegate(playerDelegate{})
			return m, nil
		}
	}

	// Forward update to the focused list
	var cmd tea.Cmd
	if m.focusedList == 0 {
		m.fieldedList, cmd = m.fieldedList.Update(msg)
	} else {
		m.benchList, cmd = m.benchList.Update(msg)
	}
	return m, cmd
}

func (m *SubstitutionModel) View() string {
	// Border styles
	focusedBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("170")).
		Padding(1)

	unfocusedBorder := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1)

	// Title styles
	focusedTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("170")).
		Align(lipgloss.Center)

	unfocusedTitle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("241")).
		Align(lipgloss.Center)

	// Left panel: Fielded players
	fieldedTitle := "Fielded Players"
	fieldedPanel := m.fieldedList.View()
	if m.focusedList == 0 {
		fieldedPanel = focusedBorder.Render(lipgloss.JoinVertical(lipgloss.Left, focusedTitle.Render(fieldedTitle), "", fieldedPanel))
	} else {
		fieldedPanel = unfocusedBorder.Render(lipgloss.JoinVertical(lipgloss.Left, unfocusedTitle.Render(fieldedTitle), "", fieldedPanel))
	}

	// Right panel: Bench players
	benchTitle := "Bench"
	benchPanel := m.benchList.View()
	if m.focusedList == 1 {
		benchPanel = focusedBorder.Render(lipgloss.JoinVertical(lipgloss.Left, focusedTitle.Render(benchTitle), "", benchPanel))
	} else {
		benchPanel = unfocusedBorder.Render(lipgloss.JoinVertical(lipgloss.Left, unfocusedTitle.Render(benchTitle), "", benchPanel))
	}

	// Join panels side by side
	panels := lipgloss.JoinHorizontal(lipgloss.Top, fieldedPanel, "  ", benchPanel)

	// Selection status - show what's been selected
	var statusText string
	if m.selectedXI != nil || m.selectedBench != nil {
		statusLines := []string{}
		if m.selectedXI != nil {
			statusLines = append(statusLines, fmt.Sprintf("Fielded: %s (%s)", m.selectedXI.Player.Name, m.selectedXI.Position))
		} else {
			statusLines = append(statusLines, "Fielded: (select player)")
		}
		if m.selectedBench != nil {
			statusLines = append(statusLines, fmt.Sprintf("Bench: %s", m.selectedBench.Player.Name))
		} else {
			statusLines = append(statusLines, "Bench: (select player)")
		}

		statusText = strings.Join(statusLines, " ⇄ ")
		if m.selectedXI != nil && m.selectedBench != nil {
			statusText = "✓ Ready: " + statusText
		}

		status := "\n" + lipgloss.NewStyle().
			Foreground(lipgloss.Color("170")).
			Bold(true).
			Render(statusText)

		return lipgloss.JoinVertical(lipgloss.Center, panels, status)
	}

	return panels
}

// playerItem wraps a MatchPlayerParticipant for use in bubbles/list
type playerItem struct {
	player  *domain.MatchPlayerParticipant
	isBench bool
}

func (i playerItem) FilterValue() string {
	return i.player.Player.Name
}

// playerDelegate handles rendering of player items
type playerDelegate struct {
	selectedPlayer *domain.MatchPlayerParticipant
}

func (d playerDelegate) Height() int                             { return 1 }
func (d playerDelegate) Spacing() int                            { return 0 }
func (d playerDelegate) Update(_ tea.Msg, _ *list.Model) tea.Cmd { return nil }
func (d playerDelegate) Render(w io.Writer, m list.Model, index int, listItem list.Item) {
	item, ok := listItem.(playerItem)
	if !ok {
		return
	}

	player := item.player

	// Format stamina as visual bars
	staminaBars := renderStaminaBars(int(player.Stamina))

	// Format the player info
	var str string
	if item.isBench {
		// Bench players don't have positions
		str = fmt.Sprintf("%-20s ★%-2d %s", player.Player.Name, player.Player.Quality, staminaBars)
	} else {
		// Current XI shows position
		str = fmt.Sprintf("%-3s %-16s ★%-2d %s", player.Position, player.Player.Name, player.Player.Quality, staminaBars)
	}

	// Determine prefix and style based on selection state
	prefix := "  "
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))

	// Check if this player is selected (compare by Player object, not pointer)
	if d.selectedPlayer != nil && player.Player == d.selectedPlayer.Player {
		prefix = "✓ "
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Bold(true)
	} else if index == m.Index() {
		// Cursor on this item
		prefix = "▸ "
		style = lipgloss.NewStyle().Foreground(lipgloss.Color("170")).Bold(true)
	}

	fmt.Fprint(w, style.Render(prefix+str))
}

// renderStaminaBars creates a visual representation of stamina
func renderStaminaBars(stamina int) string {
	filled := stamina / 20 // Each bar represents 20%
	empty := 5 - filled

	bars := strings.Repeat("●", filled) + strings.Repeat("○", empty)
	return fmt.Sprintf("%s %d%%", bars, stamina)
}

// Messages for substitution flow
type executeSubstitutionMsg struct {
	playerOut *domain.MatchPlayerParticipant
	playerIn  *domain.MatchPlayerParticipant
}
