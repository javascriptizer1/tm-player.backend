package app

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-player.backend/internal/service"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

type CityRepository interface {
	Upsert(ctx context.Context, city *domain.City) (uuid.UUID, error)
	One(ctx context.Context, cityID uuid.UUID) (*domain.City, error)
	CheckOneByName(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, args pagination.PaginationParams) ([]*domain.City, error)
	Count(ctx context.Context) (int64, error)
}

type PositionRepository interface {
	Upsert(ctx context.Context, position *domain.Position) (uuid.UUID, error)
	One(ctx context.Context, id uuid.UUID) (*domain.Position, error)
	CheckOneByName(ctx context.Context, name string) (bool, error)
	List(ctx context.Context, args service.ManyPositionsListOptions) ([]*domain.Position, error)
	Count(ctx context.Context, args service.ManyPositionsListOptions) (int64, error)
}

type PlayerRepository interface {
	Upsert(ctx context.Context, position *domain.Player) (uuid.UUID, error)
	One(ctx context.Context, id uuid.UUID) (*domain.Player, error)
	List(ctx context.Context, args service.ManyPlayersListOptions) ([]*domain.Player, error)
	Count(ctx context.Context, args service.ManyPlayersListOptions) (int64, error)
}

type CityService interface {
	Create(ctx context.Context, name string) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, dto service.CityUpdateValues) (uuid.UUID, error)
	One(ctx context.Context, id uuid.UUID) (*domain.City, error)
	List(ctx context.Context, args pagination.PaginationParams) ([]*domain.City, int64, error)
}

type PositionService interface {
	Create(ctx context.Context, name string) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, dto service.PositionUpdateValues) (uuid.UUID, error)
	One(ctx context.Context, id uuid.UUID) (*domain.Position, error)
	List(ctx context.Context, args service.ManyPositionsListOptions) ([]*domain.Position, int64, error)
}

type PlayerService interface {
	Create(ctx context.Context, dto service.PlayerCreateValues) (uuid.UUID, error)
	Update(ctx context.Context, id uuid.UUID, dto service.PlayerUpdateValues) (uuid.UUID, error)
	One(ctx context.Context, id uuid.UUID) (*domain.Player, error)
	List(ctx context.Context, args service.ManyPlayersListOptions) ([]*domain.Player, int64, error)
}
