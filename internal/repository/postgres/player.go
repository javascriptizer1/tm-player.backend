package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/javascriptizer1/tm-player.backend/internal/converter"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-player.backend/internal/domain/vo"
	"github.com/javascriptizer1/tm-player.backend/internal/service"
	"github.com/javascriptizer1/tm-player.backend/pkg/gensqlc"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
)

type playerRepository struct {
	pool    *pgxpool.Pool
	queries *gensqlc.Queries
}

func NewPlayerRepository(pool *pgxpool.Pool, queries *gensqlc.Queries) *playerRepository {
	return &playerRepository{
		pool:    pool,
		queries: queries,
	}
}

func (r *playerRepository) Upsert(ctx context.Context, player *domain.Player) (uuid.UUID, error) {
	tx, err := r.pool.Begin(ctx)
	if err != nil {
		return uuid.Nil, err
	}
	defer tx.Rollback(ctx)

	qtx := r.queries.WithTx(tx)

	playerID, err := qtx.UpsertPlayer(ctx, gensqlc.UpsertPlayerParams{
		ID:          player.ID(),
		FirstName:   player.FirstName(),
		LastName:    player.LastName(),
		MiddleName:  player.MiddleName(),
		Birthday:    player.Birthday(),
		Photo:       player.Photo(),
		CityID:      player.CityID(),
		Height:      int32(player.Height()),
		ImpactLeg:   string(player.ImpactLeg()),
		MarketValue: int32(player.MarketValue()),
	})
	if err != nil {
		return uuid.Nil, err
	}

	_, err = qtx.TrimNotExistingPlayerPositions(ctx, gensqlc.TrimNotExistingPlayerPositionsParams{
		PlayerID:            playerID,
		ExistingPositionIds: extractPositionIDs(player.Positions()),
	})
	if err != nil {
		return uuid.Nil, err
	}

	err = qtx.UpsertPlayerPositions(ctx, gensqlc.UpsertPlayerPositionsParams{
		PlayerID:    playerID,
		PositionIds: extractPositionIDs(player.Positions()),
		Mains:       extractPositionMains(player.Positions()),
	})
	if err != nil {
		return uuid.Nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return uuid.Nil, err
	}

	return playerID, nil
}

func (r *playerRepository) One(ctx context.Context, id uuid.UUID) (*domain.Player, error) {
	player, err := r.queries.GetPlayerByID(ctx, gensqlc.GetPlayerByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerror.NotFound.New("PlayerNotFound", fmt.Sprintf("player with id = %s not found", id.String()))
		}
		return nil, err
	}

	playerPositions, err := r.queries.GetPositionsForPlayers(ctx, gensqlc.GetPositionsForPlayersParams{PlayerIds: []uuid.UUID{player.ID}})
	if err != nil {
		return nil, err
	}

	result, err := converter.PlayerFromRepoToDomain(player, playerPositions)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r *playerRepository) List(ctx context.Context, args service.ManyPlayersListOptions) ([]*domain.Player, error) {
	players, err := r.queries.ListPlayers(ctx, gensqlc.ListPlayersParams{
		CityID:     args.CityID,
		PositionID: args.PositionID,
		MaxAge:     args.MaxAge,
		MinAge:     args.MinAge,
		Ids:        args.IDs,
		Limit:      args.Limit,
		Offset:     args.Offset,
	})
	if err != nil {
		return nil, err
	}

	playerIDs := make([]uuid.UUID, len(players))
	for i, player := range players {
		playerIDs[i] = player.ID
	}

	positions, err := r.queries.GetPositionsForPlayers(ctx, gensqlc.GetPositionsForPlayersParams{PlayerIds: playerIDs})
	if err != nil {
		return nil, err
	}

	positionMap := make(map[uuid.UUID][]*gensqlc.PlayerPosition)
	for _, position := range positions {
		positionMap[position.PlayerID] = append(positionMap[position.PlayerID], position)
	}

	result := make([]*domain.Player, len(players))
	for i, repoPlayer := range players {
		positions := positionMap[repoPlayer.ID]
		player, err := converter.PlayerFromRepoToDomain(repoPlayer, positions)
		if err != nil {
			return nil, err
		}
		result[i] = player
	}

	return result, nil
}

func (r *playerRepository) Count(ctx context.Context, args service.ManyPlayersListOptions) (int64, error) {
	return r.queries.CountPlayers(ctx, gensqlc.CountPlayersParams{
		CityID:     args.CityID,
		PositionID: args.PositionID,
		MaxAge:     args.MaxAge,
		MinAge:     args.MinAge,
		Ids:        args.IDs,
	})
}

func extractPositionIDs(positions []vo.PlayerPosition) []uuid.UUID {
	positionIDs := make([]uuid.UUID, len(positions))
	for i, pos := range positions {
		positionIDs[i] = pos.PositionID()
	}
	return positionIDs
}

func extractPositionMains(positions []vo.PlayerPosition) []bool {
	mains := make([]bool, len(positions))
	for i, pos := range positions {
		mains[i] = pos.Main()
	}
	return mains
}
