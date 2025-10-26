package simulation

import (
	"testing"
	"time"

	"github.com/cameronjpr/gaffer/internal/domain"
)

func TestMatchControllerSendsInitialEvent(t *testing.T) {
	// Setup test clubs
	testDB, queries := setupTestDB(t)
	defer testDB.Close()

	homeClub, awayClub := getTestClubs(t, queries)
	fixture := &domain.Fixture{HomeTeam: homeClub, AwayTeam: awayClub}
	match := domain.NewMatchFromFixture(fixture)

	// Create controller
	controller := NewMatchController(match)

	// Start controller in goroutine
	go controller.Run()

	// Wait for initial event with timeout
	select {
	case msg := <-controller.EventChan():
		switch msg := msg.(type) {
		case MatchUpdateMsg:
			if msg.Match.CurrentMinute != 1 {
				t.Errorf("Expected minute 1, got %d", msg.Match.CurrentMinute)
			}
			t.Logf("✓ Received initial event at minute %d", msg.Match.CurrentMinute)
		default:
			t.Errorf("Expected MatchUpdateMsg, got %T", msg)
		}
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for initial event - controller may be stuck")
	}
}

func TestMatchControllerSimulatesMatch(t *testing.T) {
	// Setup test clubs
	testDB, queries := setupTestDB(t)
	defer testDB.Close()

	homeClub, awayClub := getTestClubs(t, queries)
	fixture := &domain.Fixture{HomeTeam: homeClub, AwayTeam: awayClub}
	match := domain.NewMatchFromFixture(fixture)

	// Create controller
	controller := NewMatchController(match)

	// Start controller in goroutine
	go controller.Run()

	eventsReceived := 0
	maxEvents := 10 // Only receive first 10 events

	for eventsReceived < maxEvents {
		select {
		case msg := <-controller.EventChan():
			eventsReceived++
			switch msg := msg.(type) {
			case MatchUpdateMsg:
				t.Logf("Event %d: Minute %d, Score %d-%d",
					eventsReceived,
					msg.Match.CurrentMinute,
					msg.Match.Home.Score,
					msg.Match.Away.Score)
			case HalftimeMsg:
				t.Logf("Event %d: Halftime - Score %d-%d",
					eventsReceived,
					msg.Match.Home.Score,
					msg.Match.Away.Score)
			case FulltimeMsg:
				t.Logf("Event %d: Fulltime - Score %d-%d",
					eventsReceived,
					msg.Match.Home.Score,
					msg.Match.Away.Score)
				return // Match finished
			default:
				t.Logf("Event %d: %T", eventsReceived, msg)
			}
		case <-time.After(5 * time.Second):
			t.Fatalf("Timeout after receiving %d events - expected %d", eventsReceived, maxEvents)
		}
	}

	t.Logf("✓ Successfully received %d events", eventsReceived)
}

func TestMatchControllerPauseResume(t *testing.T) {
	// Setup test clubs
	testDB, queries := setupTestDB(t)
	defer testDB.Close()

	homeClub, awayClub := getTestClubs(t, queries)
	fixture := &domain.Fixture{HomeTeam: homeClub, AwayTeam: awayClub}
	match := domain.NewMatchFromFixture(fixture)

	// Create controller
	controller := NewMatchController(match)

	// Start controller
	go controller.Run()

	// Get initial event
	<-controller.EventChan()

	// Let it run for a few events
	for i := 0; i < 3; i++ {
		<-controller.EventChan()
	}

	// Pause
	controller.SendCommand(TogglePausedCmd{})

	// Should receive pause message
	select {
	case msg := <-controller.EventChan():
		if _, ok := msg.(MatchPausedMsg); !ok {
			t.Errorf("Expected MatchPausedMsg, got %T", msg)
		}
		t.Logf("✓ Received pause confirmation")
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for pause confirmation")
	}

	// Wait a bit - should not receive any more events while paused
	select {
	case msg := <-controller.EventChan():
		t.Errorf("Received unexpected event while paused: %T", msg)
	case <-time.After(500 * time.Millisecond):
		t.Logf("✓ No events received while paused (correct)")
	}

	// Resume
	controller.SendCommand(TogglePausedCmd{})

	// Should start receiving events again
	select {
	case msg := <-controller.EventChan():
		if _, ok := msg.(MatchUpdateMsg); !ok {
			t.Errorf("Expected MatchUpdateMsg after resume, got %T", msg)
		}
		t.Logf("✓ Received event after resume")
	case <-time.After(2 * time.Second):
		t.Fatal("Timeout waiting for event after resume")
	}
}
