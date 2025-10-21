package game

import "fmt"

type CommentaryMessage struct {
	Message string
	Flash   bool
	Event   EventType
}

func getCommentaryForEvent(event EventType, participant *MatchParticipant, match *Match) CommentaryMessage {
	if match == nil {
		return CommentaryMessage{Message: "Error reading match data", Flash: false}
	}

	msg := CommentaryMessage{Event: event}

	switch event {
	case HalfStartsEvent:
		switch match.CurrentHalf {
		case 1:
			msg.Message = "First half starts!"
		case 2:
			msg.Message = "Second half starts!"
		}
	case HalfEndsEvent:
		switch match.CurrentHalf {
		case 1:
			msg.Message = fmt.Sprintf("First half ends, with the score at %d-%d", match.Home.Score, match.Away.Score)
		case 2:
			msg.Message = "Full time!"
		}
	case HomeGoalScoredEvent:
		msg.Message = fmt.Sprintf("GOAL: %s score!", match.Home.Club.Name)
	case AwayGoalScoredEvent:
		msg.Message = fmt.Sprintf("GOAL: %s score!", match.Away.Club.Name)
	case PossessionChangedEvent:
		msg.Message = fmt.Sprintf("%s win the ball", participant.Club.Name)
	case PossessionRetainedEvent:
		msg.Message = fmt.Sprintf("%s have the ball...", participant.Club.Name)
	case UnknownEvent:
		msg.Message = fmt.Sprintf("Unknown event: %d", event)
	default:
		msg.Message = fmt.Sprintf("Unknown event: %d", event)
	}

	return msg
}
