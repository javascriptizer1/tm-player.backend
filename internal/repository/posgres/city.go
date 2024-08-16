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
	"github.com/javascriptizer1/tm-player.backend/pkg/gensqlc"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

type cityRepository struct {
	pool    *pgxpool.Pool
	queries *gensqlc.Queries
}

func NewCityRepository(pool *pgxpool.Pool, queries *gensqlc.Queries) *cityRepository {
	return &cityRepository{
		pool:    pool,
		queries: queries,
	}
}

func (r *cityRepository) Upsert(ctx context.Context, city *domain.City) (uuid.UUID, error) {
	return r.queries.UpsertCity(ctx, gensqlc.UpsertCityParams{
		ID:   city.ID(),
		Name: city.Name(),
	})
}

func (r *cityRepository) One(ctx context.Context, id uuid.UUID) (*domain.City, error) {
	row, err := r.queries.GetCityByID(ctx, gensqlc.GetCityByIDParams{ID: id})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domainerror.NotFound.New("CityNotFound", fmt.Sprintf("city with id = %s not found", id.String()))

		}

		return nil, err
	}

	return converter.CityFromRepoToDomain(row)
}

func (r *cityRepository) CheckOneByName(ctx context.Context, name string) (bool, error) {
	return r.queries.GetCityByNameExist(ctx, gensqlc.GetCityByNameExistParams{Name: name})
}

func (r *cityRepository) List(ctx context.Context, args pagination.PaginationParams) ([]*domain.City, error) {
	rows, err := r.queries.ListCities(ctx, gensqlc.ListCitiesParams{
		Limit:  args.Limit,
		Offset: args.Offset,
	})
	if err != nil {
		return nil, err
	}

	cities := make([]*domain.City, len(rows))
	for i, row := range rows {
		cities[i], err = converter.CityFromRepoToDomain(row)
		if err != nil {
			return []*domain.City{}, err
		}
	}

	return cities, nil
}

func (r *cityRepository) Count(ctx context.Context) (int64, error) {
	return r.queries.CountCities(ctx)
}
