package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-player.backend/internal/service"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

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
