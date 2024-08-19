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
	"github.com/javascriptizer1/tm-player.backend/internal/service"
	"github.com/javascriptizer1/tm-player.backend/pkg/gensqlc"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
)

type positionRepository struct {
	pool    *pgxpool.Pool
	queries *gensqlc.Queries
}

func NewPositionRepository(pool *pgxpool.Pool, queries *gensqlc.Queries) *positionRepository {
	return &positionRepository{
		pool:    pool,
		queries: queries,
	}
}

func (r *positionRepository) Upsert(ctx context.Context, position *domain.Position) (uuid.UUID, error) {
	return r.queries.UpsertPosition(ctx, gensqlc.UpsertPositionParams{
		ID:   position.ID(),
		Name: position.Name(),
	})
}

func (r *positionRepository) One(ctx context.Context, id uuid.UUID) (*domain.Position, error) {
	row, err := r.queries.GetPositionByID(ctx, gensqlc.GetPositionByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerror.NotFound.New("PositionNotFound", fmt.Sprintf("position with id = %s not found", id.String()))

		}

		return nil, err
	}

	return converter.PositionFromRepoToDomain(row)
}

func (r *positionRepository) CheckOneByName(ctx context.Context, name string) (bool, error) {
	return r.queries.GetPositionByNameExist(ctx, gensqlc.GetPositionByNameExistParams{Name: name})
}

func (r *positionRepository) List(ctx context.Context, args service.ManyPositionsListOptions) ([]*domain.Position, error) {
	options := gensqlc.ListPositionsParams{
		Limit:  args.Limit,
		Offset: args.Offset,
		Ids:    args.IDs,
	}

	rows, err := r.queries.ListPositions(ctx, options)
	if err != nil {
		return nil, err
	}

	positions := make([]*domain.Position, len(rows))
	for i, row := range rows {
		positions[i], err = converter.PositionFromRepoToDomain(row)
		if err != nil {
			return []*domain.Position{}, err
		}
	}

	return positions, nil
}

func (r *positionRepository) Count(ctx context.Context, args service.ManyPositionsListOptions) (int64, error) {
	return r.queries.CountPositions(ctx, gensqlc.CountPositionsParams{Ids: args.IDs})
}
