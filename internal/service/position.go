package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

type positionService struct {
	positionRepo PositionRepository
}

type PositionUpdateValues struct {
	Name string
}

type ManyPositionsListOptions struct {
	pagination.PaginationParams
	IDs []uuid.UUID
}

func NewPositionService(positionRepo PositionRepository) *positionService {
	return &positionService{positionRepo: positionRepo}
}

func (s *positionService) Create(ctx context.Context, name string) (uuid.UUID, error) {
	exists, _ := s.positionRepo.CheckOneByName(ctx, name)
	if exists {
		return uuid.Nil, domainerror.Conflict.New("PositionAlreadyExists", "position with this name already exists")
	}

	position, err := domain.NewPosition(name)
	if err != nil {
		return uuid.Nil, err
	}

	return s.positionRepo.Upsert(ctx, position)
}

func (s *positionService) Update(ctx context.Context, id uuid.UUID, dto PositionUpdateValues) (uuid.UUID, error) {
	position, err := s.positionRepo.One(ctx, id)
	if err != nil {
		return uuid.Nil, err
	}

	if position.Name() != dto.Name {
		exists, _ := s.positionRepo.CheckOneByName(ctx, dto.Name)
		if exists {
			return uuid.Nil, domainerror.Conflict.New("PositionAlreadyExists", "position with this name already exists")
		}
	}

	if err := position.SetName(dto.Name); err != nil {
		return uuid.Nil, err
	}

	return s.positionRepo.Upsert(ctx, position)
}

func (s *positionService) One(ctx context.Context, id uuid.UUID) (*domain.Position, error) {
	return s.positionRepo.One(ctx, id)
}

func (s *positionService) List(ctx context.Context, args ManyPositionsListOptions) ([]*domain.Position, int64, error) {
	var zero []*domain.Position

	positions, err := s.positionRepo.List(ctx, args)
	if err != nil {
		return zero, 0, err
	}

	count, err := s.positionRepo.Count(ctx, args)
	if err != nil {
		return zero, 0, err
	}

	return positions, count, nil
}
