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

func (g *GRPCServer) CreatePosition(ctx context.Context, input *gengrpc.CreatePositionRequest) (*gengrpc.CreatePositionResponse, error) {
	id, err := g.positionService.Create(ctx, input.GetName())
	if err != nil {
		return nil, err
	}

	return &gengrpc.CreatePositionResponse{Id: id.String()}, nil
}

func (g *GRPCServer) UpdatePosition(ctx context.Context, input *gengrpc.UpdatePositionRequest) (*gengrpc.UpdatePositionResponse, error) {
	id, err := uuid.Parse(input.GetId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
	}

	_, err = g.positionService.Update(ctx, id, service.PositionUpdateValues{Name: input.GetName()})
	if err != nil {
		return nil, err
	}

	return &gengrpc.UpdatePositionResponse{Id: id.String()}, nil
}

func (g *GRPCServer) GetPosition(ctx context.Context, input *gengrpc.GetPositionRequest) (*gengrpc.GetPositionResponse, error) {
	id, err := uuid.Parse(input.GetId())
	if err != nil {
		return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
	}

	position, err := g.positionService.One(ctx, id)
	if err != nil {
		return nil, err
	}

	result, err := converter.PositionFromDomainToProto(position)
	if err != nil {
		return nil, err
	}

	return &gengrpc.GetPositionResponse{Position: result}, nil
}

func (g *GRPCServer) GetListPositions(ctx context.Context, input *gengrpc.GetListPositionsRequest) (*gengrpc.GetListPositionsResponse, error) {
	paging, err := pagination.New(&input.Limit, &input.Page)
	if err != nil {
		return nil, err
	}

	ids := make([]uuid.UUID, len(input.GetIds()))
	for i, v := range input.GetIds() {
		ids[i], err = uuid.Parse(v)
		if err != nil {
			return nil, domainerror.Invalid.New("InvalidUUID", "invalid UUID format")
		}
	}

	positions, total, err := g.positionService.List(ctx, service.ManyPositionsListOptions{
		PaginationParams: paging,
		IDs:              ids,
	})
	if err != nil {
		return nil, err
	}

	result := make([]*gengrpc.Position, len(positions))
	for i, v := range positions {
		result[i], err = converter.PositionFromDomainToProto(v)
		if err != nil {
			return nil, err
		}
	}

	return &gengrpc.GetListPositionsResponse{Positions: result, Total: total}, nil
}
