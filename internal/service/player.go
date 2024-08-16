package service

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-player.backend/internal/domain/vo"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

type playerService struct {
	playerRepo PlayerRepository
}

type PlayerCreateValues struct {
	FirstName   string
	LastName    string
	MiddleName  *string
	Birthday    time.Time
	Photo       *string
	CityID      uuid.UUID
	Height      int64
	ImpactLeg   domain.ImpactLeg
	MarketValue int64
	Positions   []vo.PlayerPosition
}

type PlayerUpdateValues struct {
	FirstName   *string
	LastName    *string
	MiddleName  *string
	Birthday    *time.Time
	Photo       *string
	CityID      *uuid.UUID
	Height      *int64
	ImpactLeg   *domain.ImpactLeg
	MarketValue *int64
	Positions   *[]vo.PlayerPosition
}

type ManyPlayersListOptions struct {
	pagination.PaginationParams
	CityID     *uuid.UUID
	PositionID *uuid.UUID
	MinAge     *int64
	MaxAge     *int64
	IDs        []uuid.UUID
}

func NewPlayerService(playerRepo PlayerRepository) *playerService {
	return &playerService{playerRepo: playerRepo}
}

func (s *playerService) Create(ctx context.Context, dto PlayerCreateValues) (uuid.UUID, error) {
	player, err := domain.NewPlayer(
		dto.FirstName,
		dto.LastName,
		dto.MiddleName,
		dto.Birthday,
		dto.Photo,
		dto.CityID,
		dto.Positions,
		dto.Height,
		dto.ImpactLeg,
		dto.MarketValue,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return s.playerRepo.Upsert(ctx, player)
}

func (s *playerService) Update(ctx context.Context, id uuid.UUID, dto PlayerUpdateValues) (uuid.UUID, error) {
	player, err := s.playerRepo.One(ctx, id)
	if err != nil {
		return uuid.Nil, err
	}

	err = player.UpdateDetails(
		dto.FirstName,
		dto.LastName,
		dto.MiddleName,
		dto.Birthday,
		dto.Photo,
		dto.CityID,
		dto.Height,
		dto.ImpactLeg,
		dto.MarketValue,
		dto.Positions,
	)
	if err != nil {
		return uuid.Nil, err
	}

	return s.playerRepo.Upsert(ctx, player)
}

func (s *playerService) One(ctx context.Context, id uuid.UUID) (*domain.Player, error) {
	return s.playerRepo.One(ctx, id)
}

func (s *playerService) List(ctx context.Context, args ManyPlayersListOptions) ([]*domain.Player, int64, error) {
	var zero []*domain.Player

	players, err := s.playerRepo.List(ctx, args)
	if err != nil {
		return zero, 0, err
	}

	count, err := s.playerRepo.Count(ctx, args)
	if err != nil {
		return zero, 0, err
	}

	return players, count, nil
}
