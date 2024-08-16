package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

type cityService struct {
	cityRepo CityRepository
}

type CityUpdateValues struct {
	Name string
}

func NewCityService(cityRepo CityRepository) *cityService {
	return &cityService{cityRepo: cityRepo}
}

func (s *cityService) Create(ctx context.Context, name string) (uuid.UUID, error) {
	exists, _ := s.cityRepo.CheckOneByName(ctx, name)
	if exists {
		return uuid.Nil, domainerror.Conflict.New("CityAlreadyExists", "city with this name already exists")
	}

	city, err := domain.NewCity(name)
	if err != nil {
		return uuid.Nil, err
	}

	return s.cityRepo.Upsert(ctx, city)
}

func (s *cityService) Update(ctx context.Context, id uuid.UUID, dto CityUpdateValues) (uuid.UUID, error) {
	city, err := s.cityRepo.One(ctx, id)
	if err != nil {
		return uuid.Nil, err
	}

	if city.Name() != dto.Name {
		exists, _ := s.cityRepo.CheckOneByName(ctx, dto.Name)
		if exists {
			return uuid.Nil, domainerror.Conflict.New("CityAlreadyExists", "city with this name already exists")
		}
	}

	if err := city.SetName(dto.Name); err != nil {
		return uuid.Nil, err
	}

	return s.cityRepo.Upsert(ctx, city)
}

func (s *cityService) One(ctx context.Context, id uuid.UUID) (*domain.City, error) {
	return s.cityRepo.One(ctx, id)
}

func (s *cityService) List(ctx context.Context, args pagination.PaginationParams) ([]*domain.City, int64, error) {
	var zero []*domain.City

	cities, err := s.cityRepo.List(ctx, args)
	if err != nil {
		return zero, 0, err
	}

	count, err := s.cityRepo.Count(ctx)
	if err != nil {
		return zero, 0, err
	}

	return cities, count, nil
}
