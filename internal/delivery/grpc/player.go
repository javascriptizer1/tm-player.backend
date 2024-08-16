package grpc

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/converter"
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-player.backend/internal/domain/vo"
	"github.com/javascriptizer1/tm-player.backend/internal/service"
	"github.com/javascriptizer1/tm-player.backend/pkg/gengrpc"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
	"github.com/javascriptizer1/tm-shared.backend/pagination"
)

func (g *GRPCServer) CreatePlayer(ctx context.Context, input *gengrpc.CreatePlayerRequest) (*gengrpc.CreatePlayerResponse, error) {
	cityID, err := uuid.Parse(input.GetCityId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidCityUUID", "invalid City UUID format")
	}

	positions := make([]vo.PlayerPosition, len(input.GetPositions()))
	for i, pos := range input.GetPositions() {
		id, err := uuid.Parse(pos.GetId())
		if err != nil {
			return nil, domainerror.Invalid.New("InvalidPositionUUID", "invalid Position UUID format")
		}

		positions[i] = vo.NewPlayerPosition(id, pos.GetMain())
	}

	id, err := g.playerService.Create(ctx, service.PlayerCreateValues{
		FirstName:   input.GetFirstName(),
		LastName:    input.GetLastName(),
		MiddleName:  &input.MiddleName,
		Birthday:    input.GetBirthday().AsTime(),
		Photo:       &input.Photo,
		CityID:      cityID,
		Height:      input.GetHeight(),
		ImpactLeg:   converter.ProtoImpactLegDomain[input.GetImpactLeg()],
		MarketValue: 0,
		Positions:   positions,
	})
	if err != nil {
		return nil, err
	}

	return &gengrpc.CreatePlayerResponse{Id: id.String()}, nil
}

func (g *GRPCServer) UpdatePlayer(ctx context.Context, input *gengrpc.UpdatePlayerRequest) (*gengrpc.UpdatePlayerResponse, error) {
	id, err := uuid.Parse(input.GetId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
	}

	var (
		birthday  *time.Time
		cityID    *uuid.UUID
		impactLeg *domain.ImpactLeg
		positions []vo.PlayerPosition
	)

	if input.GetCityId() != "" {
		c, err := uuid.Parse(input.GetCityId())
		if err != nil {
			return nil, domainerror.Invalid.New("InvalidCityUUID", "invalid City UUID format")
		}
		cityID = &c
	}

	if input.GetBirthday() != nil {
		b := input.Birthday.AsTime()
		birthday = &b
	}

	if input.GetImpactLeg() != gengrpc.ImpactLeg_UNKNOWN {
		il := converter.ProtoImpactLegDomain[input.ImpactLeg]
		impactLeg = &il
	}

	if len(input.GetPositions()) > 0 {
		for _, pos := range input.GetPositions() {
			id, err := uuid.Parse(pos.GetId())
			if err != nil {
				return nil, domainerror.Invalid.New("InvalidPositionUUID", "invalid Position UUID format")
			}

			positions = append(positions, vo.NewPlayerPosition(id, pos.GetMain()))
		}
	}

	_, err = g.playerService.Update(ctx, id, service.PlayerUpdateValues{
		FirstName:   &input.FirstName,
		LastName:    &input.LastName,
		MiddleName:  &input.MiddleName,
		Birthday:    birthday,
		Photo:       &input.Photo,
		CityID:      cityID,
		Height:      &input.Height,
		MarketValue: nil,
		ImpactLeg:   impactLeg,
		Positions:   &positions,
	})
	if err != nil {
		return nil, err
	}

	return &gengrpc.UpdatePlayerResponse{Id: id.String()}, nil
}

func (g *GRPCServer) GetPlayer(ctx context.Context, input *gengrpc.GetPlayerRequest) (*gengrpc.GetPlayerResponse, error) {
	id, err := uuid.Parse(input.GetId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
	}

	player, err := g.playerService.One(ctx, id)
	if err != nil {
		return nil, err
	}

	result := converter.PlayerFromDomainToProto(player)

	return &gengrpc.GetPlayerResponse{Player: result}, nil
}

func (g *GRPCServer) GetListPlayers(ctx context.Context, input *gengrpc.GetListPlayersRequest) (*gengrpc.GetListPlayersResponse, error) {
	paging, err := pagination.New(&input.Limit, &input.Page)
	if err != nil {
		return nil, err
	}

	var cityID, positionID *uuid.UUID
	if input.CityId != "" {
		cityUUID, err := uuid.Parse(input.CityId)
		if err != nil {
			return nil, domainerror.Invalid.New("InvalidCityUUID", "invalid City UUID format")
		}
		cityID = &cityUUID
	}

	if input.PositionId != "" {
		positionUUID, err := uuid.Parse(input.PositionId)
		if err != nil {
			return nil, domainerror.Invalid.New("InvalidPositionUUID", "invalid Position UUID format")
		}
		positionID = &positionUUID
	}

	ids := make([]uuid.UUID, len(input.Ids))
	for i, v := range input.Ids {
		ids[i], err = uuid.Parse(v)
		if err != nil {
			return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
		}
	}

	var minAge, maxAge *int64
	if input.MinAge > 0 {
		minAge = &input.MinAge
	}
	if input.MaxAge > 0 {
		maxAge = &input.MaxAge
	}

	options := service.ManyPlayersListOptions{
		PaginationParams: paging,
		CityID:           cityID,
		PositionID:       positionID,
		MinAge:           minAge,
		MaxAge:           maxAge,
		IDs:              ids,
	}

	players, total, err := g.playerService.List(ctx, options)
	if err != nil {
		return nil, err
	}

	result := make([]*gengrpc.Player, len(players))
	for i, v := range players {
		result[i] = converter.PlayerFromDomainToProto(v)
	}

	return &gengrpc.GetListPlayersResponse{Players: result, Total: total}, nil
}
