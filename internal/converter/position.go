package converter

import (
	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-player.backend/pkg/gengrpc"
	"github.com/javascriptizer1/tm-player.backend/pkg/gensqlc"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func PositionFromRepoToDomain(repoPosition *gensqlc.Position) (*domain.Position, error) {
	if repoPosition == nil {
		return nil, domainerror.Invalid.New("NilRepoPosition", "repoPosition is nil")
	}

	return domain.NewPositionWithID(
		repoPosition.ID,
		repoPosition.Name,
		repoPosition.CreatedAt,
		repoPosition.UpdatedAt,
	)
}

func PositionFromDomainToProto(position *domain.Position) (*gengrpc.Position, error) {
	if position == nil {
		return nil, domainerror.Invalid.New("NilDomainPosition", "position is nil")
	}

	return &gengrpc.Position{
		Id:        position.ID().String(),
		Name:      position.Name(),
		CreatedAt: timestamppb.New(position.CreatedAt()),
		UpdatedAt: timestamppb.New(position.UpdatedAt()),
	}, nil
}

func PositionFromProtoToDomain(protoPosition *gengrpc.Position) (*domain.Position, error) {
	if protoPosition == nil {
		return nil, domainerror.Invalid.New("NilProtoPosition", "protoPosition is nil")
	}

	id, err := uuid.Parse(protoPosition.GetId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
	}

	return domain.NewPositionWithID(
		id,
		protoPosition.GetName(),
		protoPosition.GetCreatedAt().AsTime(),
		protoPosition.GetUpdatedAt().AsTime(),
	)
}
