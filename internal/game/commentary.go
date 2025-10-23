package game

import "fmt"

// CommentaryLine represents a line of commentary for display
type CommentaryLine struct {
	Message   string
	For       *MatchParticipant
	EventType EventType
}

// GenerateCommentary creates commentary from an event with full match context
func GenerateCommentary(event Event, match *Match) CommentaryLine {
	if match == nil || event.For == nil {
		return CommentaryLine{Message: "Error reading match data", EventType: event.Type}
	}

	line := CommentaryLine{
		For:       event.For,
		EventType: event.Type,
	}

	switch event.Type {
	case HalfStartsEvent:
		switch match.CurrentHalf {
		case 1:
			line.Message = "First half starts!"
		case 2:
			line.Message = "Second half starts!"
		}
	case HalfEndsEvent:
		switch match.CurrentHalf {
		case 1:
			line.Message = fmt.Sprintf("First half ends, with the score at %d-%d", match.Home.Score, match.Away.Score)
		case 2:
			line.Message = "Full time!"
		}
	case GoalEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("GOAL: %s scores for %s!", event.Player.Player.Name, event.For.Club.Name)
		} else {
			line.Message = fmt.Sprintf("GOAL: %s score!", event.For.Club.Name)
		}
	case PossessionChangedEvent:
		line.Message = fmt.Sprintf("%s win the ball", event.For.Club.Name)
	case PossessionRetainedEvent:
		line.Message = fmt.Sprintf("%s have the ball...", event.For.Club.Name)
	case SavedShotEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("Save by %s!", event.Player.Player.Name)
		} else {
			line.Message = "Great save!"
		}
	case MissedShotEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("Missed shot by %s!", event.Player.Player.Name)
		} else {
			line.Message = "Missed shot!"
		}
	case YellowCardEvent:
		if event.Player != nil && event.Player.Player != nil {
			line.Message = fmt.Sprintf("Yellow card for %s", event.Player.Player.Name)
		} else {
			line.Message = "Yellow card!"
		}
	case RedCardEvent:
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
