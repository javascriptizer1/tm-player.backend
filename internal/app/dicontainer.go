package app

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/javascriptizer1/tm-player.backend/internal/config"
	"github.com/javascriptizer1/tm-player.backend/internal/delivery/grpc"
	repository "github.com/javascriptizer1/tm-player.backend/internal/repository/postgres"
	"github.com/javascriptizer1/tm-player.backend/internal/service"
	"github.com/javascriptizer1/tm-player.backend/pkg/gensqlc"
	"github.com/javascriptizer1/tm-shared.backend/closer"
	"github.com/javascriptizer1/tm-shared.backend/logger"
	"go.uber.org/zap"
)

type diContainer struct {
	config             *config.Config
	pgPool             *pgxpool.Pool
	queries            *gensqlc.Queries
	cityRepository     CityRepository
	cityService        CityService
	positionRepository PositionRepository
	positionService    PositionService
	playerRepository   PlayerRepository
	playerService      PlayerService
	grpcServer         *grpc.GRPCServer
}

func newDIContainer() *diContainer {
	return &diContainer{}
}

func (d *diContainer) Config() *config.Config {
	if d.config == nil {
		d.config = config.MustLoad()
	}
	return d.config
}

func (d *diContainer) PostgresPool(ctx context.Context) (*pgxpool.Pool, error) {
	if d.pgPool == nil {
		config, err := pgxpool.ParseConfig(d.Config().DB.DSN())
		if err != nil {
			logger.Fatal("failed to parse database config: ", zap.String("err", err.Error()))

			return nil, err
		}

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			logger.Fatal("failed to create database pool: ", zap.String("err", err.Error()))

			return nil, err
		}

		closer.Add(func() error {
			pool.Close()
			return nil
		})

		d.pgPool = pool
	}

	return d.pgPool, nil
}

func (d *diContainer) Queries(ctx context.Context) (*gensqlc.Queries, error) {
	if d.queries == nil {
		db, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}

		d.queries = gensqlc.New(db)
	}

	return d.queries, nil
}

func (d *diContainer) CityRepository(ctx context.Context) (CityRepository, error) {
	if d.cityRepository == nil {
		pool, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}

		queries, err := d.Queries(ctx)
		if err != nil {
			return nil, err
		}

		d.cityRepository = repository.NewCityRepository(pool, queries)
	}

	return d.cityRepository, nil
}

func (d *diContainer) CityService(ctx context.Context) (CityService, error) {
	if d.cityService == nil {
		repo, err := d.CityRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.cityService = service.NewCityService(repo)
	}

	return d.cityService, nil
}

func (d *diContainer) PositionRepository(ctx context.Context) (PositionRepository, error) {
	if d.positionRepository == nil {
		pool, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}

		queries, err := d.Queries(ctx)
		if err != nil {
			return nil, err
		}

		d.positionRepository = repository.NewPositionRepository(pool, queries)
	}

	return d.positionRepository, nil
}

func (d *diContainer) PositionService(ctx context.Context) (PositionService, error) {
	if d.positionService == nil {
		repo, err := d.PositionRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.positionService = service.NewPositionService(repo)
	}

	return d.positionService, nil
}

func (d *diContainer) PlayerRepository(ctx context.Context) (PlayerRepository, error) {
	if d.playerRepository == nil {
		pool, err := d.PostgresPool(ctx)
		if err != nil {
			return nil, err
		}

		queries, err := d.Queries(ctx)
		if err != nil {
			return nil, err
		}

		d.playerRepository = repository.NewPlayerRepository(pool, queries)
	}

	return d.playerRepository, nil
}

func (d *diContainer) PlayerService(ctx context.Context) (PlayerService, error) {
	if d.playerService == nil {
		repo, err := d.PlayerRepository(ctx)
		if err != nil {
			return nil, err
		}

		d.playerService = service.NewPlayerService(repo)
	}

	return d.playerService, nil
}

func (d *diContainer) GRPCServer(ctx context.Context) (*grpc.GRPCServer, error) {
	if d.grpcServer == nil {
		cityService, err := d.CityService(ctx)
		if err != nil {
			return nil, err
		}

		positionService, err := d.PositionService(ctx)
		if err != nil {
			return nil, err
		}

		playerService, err := d.PlayerService(ctx)
		if err != nil {
			return nil, err
		}

		d.grpcServer = grpc.NewGRPCServer(cityService, positionService, playerService)
	}

	return d.grpcServer, nil
}
