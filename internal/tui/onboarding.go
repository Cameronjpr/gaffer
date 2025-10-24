package tui

import (
	"fmt"
	"os"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/huh"
	"github.com/charmbracelet/lipgloss"
)

type menuKeyMap struct {
	Exit     key.Binding
	formKeys *huh.KeyMap
}

func defaultMenuKeyMap() *menuKeyMap {
	keys := &menuKeyMap{
		Exit:     key.NewBinding(key.WithKeys("esc"), key.WithHelp("escape", "exit")),
		formKeys: huh.NewDefaultKeyMap(),
	}
	keys.formKeys.Quit.SetEnabled(false)
	return keys
}

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

func NewOnboardingModel(season *domain.Season) OnboardingModel {
	formData := &OnboardingFormData{}
	keys := defaultMenuKeyMap()

	// Build club options from season
	clubOptions := make([]huh.Option[string], len(season.Clubs))
	for i, club := range season.Clubs {
		clubOptions[i] = huh.NewOption(club.Name, club.Name)
	}

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Value(&formData.ClubName).
				Title("Choose your club:").
				Options(clubOptions...),
		),
	).WithWidth(60).WithKeyMap(keys.formKeys)

	return OnboardingModel{
		formData: formData,
		form:     form,
		keys:     keys,
	}
}

func (m OnboardingModel) Init() tea.Cmd {
	return m.form.Init()
}

func (m OnboardingModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (m OnboardingModel) View() string {
	output := lipgloss.JoinVertical(lipgloss.Center,
		m.form.View(),
	)
	return lipgloss.Place(m.width, m.height, lipgloss.Center, lipgloss.Center, output)
}

func main() {
	_, err := tea.NewProgram(NewModel()).Run()
	if err != nil {
		fmt.Println("Oh no:", err)
		os.Exit(1)
	}
}
