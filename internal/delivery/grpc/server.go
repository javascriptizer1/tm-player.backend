package grpc

import "github.com/javascriptizer1/tm-player.backend/pkg/gengrpc"

type GRPCServer struct {
	gengrpc.UnimplementedCityServiceServer
	gengrpc.UnimplementedPositionServiceServer
	gengrpc.UnimplementedPlayerServiceServer

	cityService     CityService
	positionService PositionService
	playerService   PlayerService
}

func NewGRPCServer(
	cityService CityService,
	positionService PositionService,
	playerService PlayerService,
) *GRPCServer {
	return &GRPCServer{
		cityService:     cityService,
		positionService: positionService,
		playerService:   playerService,
	}
}
