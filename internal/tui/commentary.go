package tui

import (
	"fmt"

	"github.com/cameronjpr/gaffer/internal/domain"
)

// CommentaryLine represents a line of commentary for display
type CommentaryLine struct {
	Message   string
	For       *domain.MatchParticipant
	EventType domain.EventType
}

// GenerateCommentary creates commentary from an event with full match context
func GenerateCommentary(event domain.Event, match *domain.Match) CommentaryLine {
	if match == nil || event.For == nil {
		return CommentaryLine{Message: "Error reading match data", EventType: event.Type}
	}

	line := CommentaryLine{
		For:       event.For,
		EventType: event.Type,
	}

	switch event.Type {
	case domain.HalfStartsEvent:
		switch match.CurrentHalf {
		case 1:
			line.Message = "First half starts!"
		case 2:
			line.Message = "Second half starts!"
		}
	case domain.HalfEndsEvent:
		switch match.CurrentHalf {
		case 1:
			line.Message = fmt.Sprintf("First half ends, with the score at %d-%d", match.Home.Score, match.Away.Score)
		case 2:
			line.Message = "Full time!"
		}
	case domain.GoalEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("GOAL: %s scores for %s!", event.Player.Player.Name, event.For.Club.Name)
		} else {
			line.Message = fmt.Sprintf("GOAL: %s score!", event.For.Club.Name)
		}
	case domain.PossessionChangedEvent:
		line.Message = fmt.Sprintf("%s win the ball", event.For.Club.Name)
	case domain.PossessionRetainedEvent:
		line.Message = fmt.Sprintf("%s have the ball...", event.For.Club.Name)
	case domain.SavedShotEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("Save by %s!", event.Player.Player.Name)
		} else {
			line.Message = "Great save!"
		}
	case domain.MissedShotEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("Missed shot by %s!", event.Player.Player.Name)
		} else {
			line.Message = "Missed shot!"
		}
	case domain.YellowCardEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("Yellow card for %s", event.Player.Player.Name)
		} else {
			line.Message = "Yellow card!"
		}
	case domain.RedCardEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("RED CARD! %s is sent off!", event.Player.Player.Name)
		} else {
			line.Message = "RED CARD!"
		}
	default:
		line.Message = fmt.Sprintf("Event: %d", event.Type)
	}

	return line
}
