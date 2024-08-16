package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

type CityRepository interface {
	Upsert(ctx context.Context, city *domain.City) (uuid.UUID, error)
	One(ctx context.Context, id uuid.UUID) (*domain.City, error)
	CheckOneByName(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, args pagination.PaginationParams) ([]*domain.City, error)
	Count(ctx context.Context) (int64, error)
}

type PositionRepository interface {
	Upsert(ctx context.Context, position *domain.Position) (uuid.UUID, error)
	One(ctx context.Context, id uuid.UUID) (*domain.Position, error)
	CheckOneByName(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, args ManyPositionsListOptions) ([]*domain.Position, error)
	Count(ctx context.Context, args ManyPositionsListOptions) (int64, error)
}

type PlayerRepository interface {
	Upsert(ctx context.Context, position *domain.Player) (uuid.UUID, error)
	One(ctx context.Context, id uuid.UUID) (*domain.Player, error)
	List(ctx context.Context, args ManyPlayersListOptions) ([]*domain.Player, error)
	Count(ctx context.Context, args ManyPlayersListOptions) (int64, error)
}
