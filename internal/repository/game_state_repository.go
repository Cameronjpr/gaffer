package repository

import (
	"context"
	"fmt"

	"github.com/cameronjpr/gaffer/internal/db"
	"github.com/cameronjpr/gaffer/internal/domain"
)

type GameStateRepo struct {
	queries *db.Queries
}

func NewGameStateRepository(queries *db.Queries) *GameStateRepo {
	return &GameStateRepo{queries: queries}
}

func dbGameStateToDomain(gameState db.GameState) *domain.GameState {
	return &domain.GameState{
		SelectedClubID: gameState.SelectedClubID,
		ManagerName:    gameState.ManagerName,
	}
}

func (r *GameStateRepo) GetMostRecentGameState() (*domain.GameState, error) {
	ctx := context.Background()
	gameState, err := r.queries.GetMostRecentGameState(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get game state: %w", err)
	}
	return dbGameStateToDomain(gameState), nil
}

var _ domain.GameStateRepository = (*GameStateRepo)(nil)
