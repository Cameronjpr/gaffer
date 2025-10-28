package components

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
	"github.com/charmbracelet/lipgloss"
)

func Clock(match *domain.Match) string {
	timeStr := fmt.Sprintf("(%v:00)", match.CurrentMinute)
	if match.IsHalfTime() {
		timeStr = "HT"
	} else if match.IsFullTime() {
		timeStr = "FT"
	} else if match.IsFirstHalf() && match.IsInAddedTime() {
		timeStr += fmt.Sprintf("+%v'", match.GetAddedTime(domain.FirstHalf))
	} else if match.IsSecondHalf() && match.IsInAddedTime() {
		timeStr += fmt.Sprintf("+%v'", match.GetAddedTime(domain.SecondHalf))
	}

	return lipgloss.NewStyle().
		Padding(0, 1).
		Render(timeStr)
}
