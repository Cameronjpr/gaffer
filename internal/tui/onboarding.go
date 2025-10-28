package tui

import (
	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type OnboardingModel struct {
	form     *huh.Form
	formData *OnboardingFormData
	keys     *menuKeyMap
	width    int
	height   int
}

type OnboardingFormData struct {
	ClubName string
}

func NewOnboardingModel(clubs []*domain.ClubWithPlayers) *OnboardingModel {
	formData := &OnboardingFormData{}
	keys := defaultMenuKeyMap()

	// Build club options from clubs
	clubOptions := make([]huh.Option[string], len(clubs))
	for i, club := range clubs {
		clubOptions[i] = huh.NewOption(club.Club.Name, club.Club.Name)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&formData.ClubName).
				Title("Choose your club:").
				Options(clubOptions...),
		),
	).WithWidth(60).WithKeyMap(keys.formKeys)

	return &OnboardingModel{
		formData: formData,
		form:     form,
		keys:     keys,
	}
}

func (m *OnboardingModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m *OnboardingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		formWidth := min(60, msg.Width-4)
		m.form = m.form.WithWidth(formWidth)
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Interrupt
		case "esc", "q":
			return m, tea.Quit
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
				ClubName: m.formData.ClubName,
			}
		})
	}

	return m, tea.Batch(cmds...)
}

func (m *OnboardingModel) View() string {
	output := lipgloss.JoinVertical(lipgloss.Center,
		m.form.View(),
	)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, output)
}
