package grpc

import (
	"context"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/converter"
	"github.com/javascriptizer1/tm-player.backend/internal/service"
	"github.com/javascriptizer1/tm-player.backend/pkg/gengrpc"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

func (g *GRPCServer) CreateCity(ctx context.Context, input *gengrpc.CreateCityRequest) (*gengrpc.CreateCityResponse, error) {
	id, err := g.cityService.Create(ctx, input.GetName())
	if err != nil {
		return nil, err
	}

	return &gengrpc.CreateCityResponse{Id: id.String()}, nil
}

func (g *GRPCServer) UpdateCity(ctx context.Context, input *gengrpc.UpdateCityRequest) (*gengrpc.UpdateCityResponse, error) {
	id, err := uuid.Parse(input.GetId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
	}

	_, err = g.cityService.Update(ctx, id, service.CityUpdateValues{Name: input.GetName()})
	if err != nil {
		return nil, err
	}

	return &gengrpc.UpdateCityResponse{Id: id.String()}, nil
}

func (g *GRPCServer) GetCity(ctx context.Context, input *gengrpc.GetCityRequest) (*gengrpc.GetCityResponse, error) {
	id, err := uuid.Parse(input.GetId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
	}

	city, err := g.cityService.One(ctx, id)
	if err != nil {
		return nil, err
	}

	result, err := converter.CityFromDomainToProto(city)
	if err != nil {
		return nil, err
	}

	return &gengrpc.GetCityResponse{City: result}, nil
}

func (g *GRPCServer) GetListCities(ctx context.Context, input *gengrpc.GetListCitiesRequest) (*gengrpc.GetListCitiesResponse, error) {
	paging, err := pagination.New(&input.Limit, &input.Page)
	if err != nil {
		return nil, err
	}

	cities, total, err := g.cityService.List(ctx, paging)
	if err != nil {
		return nil, err
	}

	result := make([]*gengrpc.City, len(cities))
	for i, v := range cities {
		result[i], err = converter.CityFromDomainToProto(v)
		if err != nil {
			return nil, err
		}
	}

	return &gengrpc.GetListCitiesResponse{Cities: result, Total: total}, nil
}
