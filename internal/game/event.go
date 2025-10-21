package game

import "fmt"

// EventType represents the type of event that occurred in the match
type EventType int

const (
	GoalEvent EventType = iota
	SaveEvent
	YellowCardEvent
	RedCardEvent
	InjuryEvent
	CornerEvent
	FreeKickEvent
	PossessionChangedEvent
	PossessionRetainedEvent
	HalfStartsEvent
	HalfEndsEvent
)

// Event represents a key moment in the match
type Event struct {
	Type   EventType
	Minute int
	For    *MatchParticipant            // Team the event benefits/involves
	Player *MatchPlayerParticipant      // Optional: involved player
}

// NewEvent creates a new event
func NewEvent(eventType EventType, minute int, participant *MatchParticipant, player *MatchPlayerParticipant) Event {
	return Event{
		Type:   eventType,
		Minute: minute,
		For:    participant,
		Player: player,
	}
}

// String returns a string representation of the event for timeline display
func (e Event) String() string {
	if e.Player != nil && e.Player.Player != nil {
		return fmt.Sprintf("%s (%d')", e.Player.Player.Name, e.Minute)
	}
	return fmt.Sprintf("(%d')", e.Minute)
}
