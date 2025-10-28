package simulation

import (
	"time"

	"github.com/cameronjpr/gaffer/internal/domain"
	tea "github.com/charmbracelet/bubbletea"
)

var speeds = []time.Duration{
	1 * time.Second,
	500 * time.Millisecond,
	250 * time.Millisecond,
	100 * time.Millisecond,
}

// MatchController orchestrates match simulation independently of the UI
type MatchController struct {
	match       *domain.Match
	engine      *Engine
	commandChan chan Command
	eventChan   chan tea.Msg
	paused      bool
	done        bool
	speedIndex  int
	ticker      *time.Ticker
}

// NewMatchController creates a new controller for the given match
func NewMatchController(match *domain.Match) *MatchController {
	return &MatchController{
		match:       match,
		engine:      NewEngine(match),
		commandChan: make(chan Command, 10), // Buffered to avoid blocking TUI
		eventChan:   make(chan tea.Msg, 10), // Buffered for smooth playback
		paused:      false,                  // Start unpaused, simulation begins immediately
		speedIndex:  2,
		ticker:      time.NewTicker(speeds[2]),
		done:        false,
	}
}

// Run starts the simulation loop (should be called in a goroutine)
func (mc *MatchController) Run() {
	// Send initial state immediately so TUI has something to render
	mc.eventChan <- MatchUpdateMsg{
		Match:       mc.match,
		LatestEvent: nil,
	}

	// Control simulation speed - adjust this to make matches faster/slower
	defer mc.ticker.Stop()

	for !mc.done {
		select {
		case cmd := <-mc.commandChan:
			mc.handleCommand(cmd)

		case <-mc.ticker.C:
			if mc.paused {
				continue // Skip simulation while paused
			}

			// Simulate one minute
			mc.engine.SimulateMinute()

			// Get the latest event if one occurred this phase
			var latestEvent *domain.Event
			if len(mc.match.Events) > 0 {
				latestEvent = &mc.match.Events[len(mc.match.Events)-1]
			}

			// Send update to TUI
			mc.eventChan <- MatchUpdateMsg{
				Match:       mc.match,
				LatestEvent: latestEvent,
			}

			// Check for halftime
			if mc.match.IsHalfTime() {
				mc.eventChan <- HalftimeMsg{Match: mc.match}
				mc.paused = true // Auto-pause at halftime
			}

			// Check for fulltime
			if mc.match.IsFullTime() {
				mc.eventChan <- FulltimeMsg{Match: mc.match}
				mc.done = true
				return
			}
		}
	}
}

// SendCommand sends a command to the controller (non-blocking)
func (mc *MatchController) SendCommand(cmd Command) {
	select {
	case mc.commandChan <- cmd:
	default:
		// Channel full, drop command (shouldn't happen with buffer)
	}
}

// EventChan returns the read-only event channel for receiving updates
func (mc *MatchController) EventChan() <-chan tea.Msg {
	return mc.eventChan
}

func (mc *MatchController) GetSpeed() string {
	switch mc.speedIndex {
	case 0:
		return "►"
	case 1:
		return "►►"
	case 2:
		return "►►►"
	case 3:
		return "►►►►"
	default:
		return "Unknown speed"
	}
}

// handleCommand processes commands from the TUI
func (mc *MatchController) handleCommand(cmd Command) {
	switch cmd.(type) {
	case PauseMatchCmd:
		mc.paused = true
		mc.eventChan <- MatchPausedMsg{Match: mc.match}

	case ResumeMatchCmd:
		mc.paused = false

	case SpeedUpCmd:
		if mc.speedIndex == len(speeds)-1 {
			mc.speedIndex = 0
		} else {
			mc.speedIndex++
		}
		// Always recreate ticker with new speed
		mc.ticker.Stop()
		mc.ticker = time.NewTicker(speeds[mc.speedIndex])

	case SlowDownCmd:
		if mc.speedIndex == 0 {
			mc.speedIndex = len(speeds) - 1
		} else {
			mc.speedIndex--
		}
		// Always recreate ticker with new speed
		mc.ticker.Stop()
		mc.ticker = time.NewTicker(speeds[mc.speedIndex])

	case StartMatchCmd:
		// Start is just resume with a different name
		mc.paused = false

	case TogglePausedCmd:
		mc.paused = !mc.paused
		if mc.paused {
			mc.eventChan <- MatchPausedMsg{Match: mc.match}
		} else {
			// When unpausing at halftime, start the second half
			if mc.match.IsHalfTime() {
				mc.match.StartSecondHalf()
				mc.eventChan <- MatchResumedMsg{Match: mc.match}
			}
		}
	// If unpausing, next tick will send MatchUpdateMsg automatically

	case SubstitutePlayerCmd:
		subCmd := cmd.(SubstitutePlayerCmd)
		subCmd.Participant.MakeSubstitution(subCmd.PlayerIn, subCmd.PlayerOut)
		mc.eventChan <- SubstitutionMadeMsg{
			Match:     mc.match,
			PlayerOut: subCmd.PlayerOut,
			PlayerIn:  subCmd.PlayerIn,
		}
	}
}
