package converter

import (
	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-player.backend/pkg/gengrpc"
	"github.com/javascriptizer1/tm-player.backend/pkg/gensqlc"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func CityFromRepoToDomain(repoCity *gensqlc.City) (*domain.City, error) {
	if repoCity == nil {
		return nil, domainerror.Invalid.New("NilRepoCity", "repoCity is nil")
	}

	return domain.NewCityWithID(
		repoCity.ID,
		repoCity.Name,
		repoCity.CreatedAt,
		repoCity.UpdatedAt,
	)
}

func CityFromDomainToProto(city *domain.City) (*gengrpc.City, error) {
	if city == nil {
		return nil, domainerror.Invalid.New("NilDomainCity", "city is nil")
	}

	return &gengrpc.City{
		Id:        city.ID().String(),
		Name:      city.Name(),
		CreatedAt: timestamppb.New(city.CreatedAt()),
		UpdatedAt: timestamppb.New(city.UpdatedAt()),
	}, nil
}

func CityFromProtoToDomain(protoCity *gengrpc.City) (*domain.City, error) {
	if protoCity == nil {
		return nil, domainerror.Invalid.New("NilProtoCity", "protoCity is nil")
	}

	id, err := uuid.Parse(protoCity.GetId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
	}

	return domain.NewCityWithID(
		id,
		protoCity.GetName(),
		protoCity.GetCreatedAt().AsTime(),
		protoCity.GetUpdatedAt().AsTime(),
	)
}
