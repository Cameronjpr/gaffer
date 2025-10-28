package simulation

import "github.com/cameronjpr/gaffer/internal/domain"

// Commands - sent from TUI to MatchController

type Command interface {
	isCommand()
}

type StartMatchCmd struct{}

func (StartMatchCmd) isCommand() {}

type PauseMatchCmd struct{}

func (PauseMatchCmd) isCommand() {}

type ResumeMatchCmd struct{}

func (ResumeMatchCmd) isCommand() {}

type TogglePausedCmd struct{}

func (TogglePausedCmd) isCommand() {}

type SpeedUpCmd struct{}

func (SpeedUpCmd) isCommand() {}

type SlowDownCmd struct{}

func (SlowDownCmd) isCommand() {}

type SubstitutePlayerCmd struct {
	Participant *domain.MatchParticipant
	PlayerOut   *domain.MatchPlayerParticipant
	PlayerIn    *domain.MatchPlayerParticipant
}

func (SubstitutePlayerCmd) isCommand() {}

// Events - sent from MatchController to TUI

// MatchUpdateMsg is sent every phase with current match state and latest event
type MatchUpdateMsg struct {
	Match       *domain.Match
	LatestEvent *domain.Event // nil if no event occurred this phase
}

// HalftimeMsg is sent when first half ends
type HalftimeMsg struct {
	Match *domain.Match
}

// FulltimeMsg is sent when match ends
type FulltimeMsg struct {
	Match *domain.Match
}

// MatchPausedMsg is sent when match is paused by user
type MatchPausedMsg struct {
	Match *domain.Match
}

// MatchResumedMsg is sent when match is resumed by user
type MatchResumedMsg struct {
	Match *domain.Match
}

type SubstitutionMadeMsg struct {
	Match     *domain.Match
	PlayerOut *domain.MatchPlayerParticipant
	PlayerIn  *domain.MatchPlayerParticipant
}
