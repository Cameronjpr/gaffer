package tui

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/components"
	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type PreMatchModel struct {
	form     *huh.Form
	formData *PreMatchFormData
	keys     *menuKeyMap
	match    *domain.Match
	width    int
	height   int
}

type PreMatchFormData struct {
	Formation string
}

func NewPreMatchModel(match *domain.Match) *PreMatchModel {

	formData := &PreMatchFormData{}
	keys := defaultMenuKeyMap()

	// Build club options from season
	formationOptions := []huh.Option[string]{
		{Key: "4-4-2", Value: "4-4-2"},
		{Key: "4-3-3", Value: "4-3-3"},
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&formData.Formation).
				Title("Choose your formation:").
				Options(formationOptions...),
		),
	).WithWidth(60).WithKeyMap(keys.formKeys)

	return &PreMatchModel{
		match:    match,
		form:     form,
		formData: formData,
		keys:     keys,
	}
}

func (m *PreMatchModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *PreMatchModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		formWidth := min(60, msg.Width-4)
		m.form = m.form.WithWidth(formWidth)
		return m, nil

	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		case tea.KeyEnter:
			return m, func() tea.Msg { return startMatchMsg{} }
		}
	}

	var cmds []tea.Cmd

	form, cmd := m.form.Update(msg)
	if f, ok := form.(*huh.Form); ok {
		m.form = f
		cmds = append(cmds, cmd)
	}

	if m.form.State == huh.StateCompleted {
		cmds = append(cmds, func() tea.Msg {
			return goToManagerHubMsg{
				ClubID: m.match.Home.Club.ID,
				/// ??????? makes no sense
			}
		})
	}

	return m, tea.Batch(cmds...)
}

func (m *PreMatchModel) View() string {
	// Calculate column width
	colWidth := m.width / 3

	// Home team section

	// Match info section
	matchContent := lipgloss.JoinVertical(lipgloss.Center,
		fmt.Sprintf("%s vs %s\n\nPress Enter to start", m.match.Home.Club.Name, m.match.Away.Club.Name),
		// m.form.View(),
	)

	homeTeamSheet := components.TeamSheet(m.match.Home)
	awayTeamSheet := components.TeamSheet(m.match.Away)

	// Create 3-column layout
	layout := lipgloss.JoinHorizontal(
		lipgloss.Top,
		lipgloss.Place(colWidth, m.height, lipgloss.Center, lipgloss.Center, homeTeamSheet),
		lipgloss.Place(colWidth, m.height, lipgloss.Center, lipgloss.Center, matchContent),
		lipgloss.Place(colWidth, m.height, lipgloss.Center, lipgloss.Center, awayTeamSheet),
	)

	return layout
}
