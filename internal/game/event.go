package game

import "fmt"

type EventType int

const (
	HalfStartsEvent EventType = iota
	HalfEndsEvent
	HomeGoalScoredEvent
	AwayGoalScoredEvent
	PossessionChangedEvent
	PossessionRetainedEvent
	UnknownEvent
)

type PlayerEventType int

const (
	PlayerScoredEvent PlayerEventType = iota
	PlayerInjuredEvent
	PlayerYellowCardEvent
	PlayerRedCardEvent
	PlayerUnknownEvent
)

type PlayerEvent struct {
	Type   PlayerEventType
	Player Player
	Minute int
}

type MatchEvent struct {
	Type   EventType
	Player Player
}

func NewPlayerEvent(t PlayerEventType, p Player, min int) PlayerEvent {
	return PlayerEvent{
		Type:   t,
		Minute: min,
		Player: p,
	}
}

func (pe PlayerEvent) String() string {
	switch pe.Type {
	case PlayerScoredEvent:
		return fmt.Sprintf("%s (%d')", pe.Player.Name, pe.Minute)
	case PlayerUnknownEvent:
		return fmt.Sprintf("%s did something", pe.Player.Name)
	default:
		return fmt.Sprintf("%s did something", pe.Player.Name)
	}
}

func NewMatchEvent(t EventType, p Player) MatchEvent {
	return MatchEvent{
		Type:   t,
		Player: p,
	}
}
