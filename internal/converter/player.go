package converter

import (
	"github.com/javascriptizer1/tm-player.backend/internal/domain"
	"github.com/javascriptizer1/tm-player.backend/internal/domain/vo"
	"github.com/javascriptizer1/tm-player.backend/pkg/gengrpc"
	"github.com/javascriptizer1/tm-player.backend/pkg/gensqlc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var ProtoImpactLegDomain = map[gengrpc.ImpactLeg]domain.ImpactLeg{
	gengrpc.ImpactLeg_UNKNOWN: domain.Unknown,
	gengrpc.ImpactLeg_LEFT:    domain.Left,
	gengrpc.ImpactLeg_RIGHT:   domain.Right,
	gengrpc.ImpactLeg_BOTH:    domain.Both,
}

var DomainImpactLegProto = map[domain.ImpactLeg]gengrpc.ImpactLeg{
	domain.Unknown: gengrpc.ImpactLeg_UNKNOWN,
	domain.Left:    gengrpc.ImpactLeg_LEFT,
	domain.Right:   gengrpc.ImpactLeg_RIGHT,
	domain.Both:    gengrpc.ImpactLeg_BOTH,
}

func PlayerFromRepoToDomain(repoPlayer *gensqlc.Player, repoPlayerPositions []*gensqlc.PlayerPosition) (*domain.Player, error) {
	var positions []vo.PlayerPosition
	for _, pos := range repoPlayerPositions {
		position := vo.NewPlayerPosition(pos.PositionID, pos.Main)
		positions = append(positions, position)
	}

	player, err := domain.NewPlayerWithID(
		repoPlayer.ID,
		repoPlayer.FirstName,
		repoPlayer.LastName,
		repoPlayer.MiddleName,
		repoPlayer.Birthday,
		repoPlayer.Photo,
		repoPlayer.CityID,
		positions,
		int64(repoPlayer.Height),
		domain.ImpactLeg(repoPlayer.ImpactLeg),
		int64(repoPlayer.MarketValue),
		repoPlayer.CreatedAt,
		repoPlayer.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return player, nil
}

func PlayerFromDomainToProto(player *domain.Player) *gengrpc.Player {
	var positions []*gengrpc.PlayerPosition
	for _, pos := range player.Positions() {
		positions = append(positions, &gengrpc.PlayerPosition{
			Id:   pos.PositionID().String(),
			Main: pos.Main(),
		})
	}

	var (
		middleName string
		photo      string
	)

	if player.MiddleName() != nil {
		middleName = *player.MiddleName()
	}

	if player.Photo() != nil {
		photo = *player.Photo()
	}

	return &gengrpc.Player{
		Id:          player.ID().String(),
		FirstName:   player.FirstName(),
		LastName:    player.LastName(),
		MiddleName:  middleName,
		Birthday:    timestamppb.New(player.Birthday()),
		Photo:       photo,
		CityId:      player.CityID().String(),
		Positions:   positions,
		Height:      player.Height(),
		ImpactLeg:   DomainImpactLegProto[player.ImpactLeg()],
		MarketValue: player.MarketValue(),
		CreatedAt:   timestamppb.New(player.CreatedAt()),
		UpdatedAt:   timestamppb.New(player.UpdatedAt()),
	}
}
