package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-player.backend/internal/domain/vo"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
)

type ImpactLeg string

const (
	Left    ImpactLeg = "left"
	Right   ImpactLeg = "right"
	Both    ImpactLeg = "both"
	Unknown ImpactLeg = "unknown"
)

type Player struct {
	id          uuid.UUID
	firstName   string
	lastName    string
	middleName  *string
	birthday    time.Time
	photo       *string
	cityID      uuid.UUID
	positions   []vo.PlayerPosition
	height      int64
	impactLeg   ImpactLeg
	marketValue int64
	createdAt   time.Time
	updatedAt   time.Time
}

func (p *Player) ID() uuid.UUID {
	return p.id
}

func (p *Player) FirstName() string {
	return p.firstName
}

func (p *Player) LastName() string {
	return p.lastName
}

func (p *Player) MiddleName() *string {
	return p.middleName
}

func (p *Player) Birthday() time.Time {
	return p.birthday
}

func (p *Player) Photo() *string {
	return p.photo
}

func (p *Player) CityID() uuid.UUID {
	return p.cityID
}

func (p *Player) Positions() []vo.PlayerPosition {
	return p.positions
}

func (p *Player) Height() int64 {
	return p.height
}

func (p *Player) ImpactLeg() ImpactLeg {
	return p.impactLeg
}

func (p *Player) MarketValue() int64 {
	return p.marketValue
}

func (p *Player) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Player) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Player) SetFirstName(name string) error {
	p.firstName = strings.TrimSpace(name)
	p.SetUpdatedAt()

	return p.validate()
}

func (p *Player) SetLastName(name string) error {
	p.lastName = strings.TrimSpace(name)
	p.SetUpdatedAt()

	return p.validate()
}

func (p *Player) SetMiddleName(name *string) error {
	if name == nil {
		p.middleName = name
		p.SetUpdatedAt()

		return nil
	}

	n := strings.TrimSpace(*name)
	p.middleName = &n
	p.SetUpdatedAt()

	return p.validate()
}

func (p *Player) SetBirthday(birthday time.Time) error {
	p.birthday = birthday
	p.SetUpdatedAt()

	return p.validate()
}

func (p *Player) SetPhoto(photo *string) {
	if photo == nil {
		p.photo = photo
		p.SetUpdatedAt()
		return
	}

	ph := strings.TrimSpace(*photo)
	p.photo = &ph
	p.SetUpdatedAt()
}

func (p *Player) SetCityID(cityID uuid.UUID) error {
	p.cityID = cityID
	p.SetUpdatedAt()

	return p.validate()
}

func (p *Player) SetHeight(height int64) error {
	p.height = height
	p.SetUpdatedAt()

	return p.validate()
}

func (p *Player) SetImpactLeg(impactLeg ImpactLeg) {
	p.impactLeg = impactLeg
	p.SetUpdatedAt()
}

func (p *Player) SetMarketValue(marketValue int64) error {
	p.marketValue = marketValue
	p.SetUpdatedAt()

	return p.validate()
}

func (p *Player) SetPositions(positions []vo.PlayerPosition) error {
	p.positions = positions
	p.SetUpdatedAt()

	return p.validate()
}

func (p *Player) SetUpdatedAt() {
	p.updatedAt = time.Now().UTC()
}

func (p *Player) UpdateDetails(
	firstName, lastName, middleName *string,
	birthday *time.Time,
	photo *string,
	cityID *uuid.UUID,
	height *int64,
	impactLeg *ImpactLeg,
	marketValue *int64,
	positions *[]vo.PlayerPosition,
) error {
	if firstName != nil {
		if err := p.SetFirstName(*firstName); err != nil {
			return err
		}
	}

	if lastName != nil {
		if err := p.SetLastName(*lastName); err != nil {
			return err
		}
	}

	if middleName != nil {
		if err := p.SetMiddleName(middleName); err != nil {
			return err
		}
	}

	if birthday != nil {
		if err := p.SetBirthday(*birthday); err != nil {
			return err
		}
	}

	if photo != nil {
		p.SetPhoto(photo)
	}

	if cityID != nil {
		if err := p.SetCityID(*cityID); err != nil {
			return err
		}
	}

	if height != nil {
		if err := p.SetHeight(*height); err != nil {
			return err
		}
	}

	if impactLeg != nil {
		p.SetImpactLeg(*impactLeg)
	}

	if marketValue != nil {
		if err := p.SetMarketValue(*marketValue); err != nil {
			return err
		}
	}

	if positions != nil {
		if err := p.SetPositions(*positions); err != nil {
			return err
		}
	}

	return nil
}

func (p *Player) validate() error {
	if len([]rune(p.firstName)) == 0 || len([]rune(p.lastName)) == 0 {
		return domainerror.Invalid.New("PlayerNameInvalid", "first or last name is invalid")
	}

	if p.height <= 0 {
		return domainerror.Invalid.New("PlayerHeightInvalid", "player height must be positive")
	}

	if len(p.positions) == 0 {
		return domainerror.Invalid.New("NoPositionsProvided", "at least one position must be provided")
	}

	if p.createdAt.After(p.updatedAt) {
		return domainerror.Invalid.New("InvalidTimeRange", "created_at must be before updated_at")
	}

	positionSet := make(map[uuid.UUID]struct{})
	for _, pos := range p.positions {
		_, duplicate := positionSet[pos.PositionID()]
		if duplicate {
			return domainerror.Invalid.New("DuplicatePosition", "duplicate positions are not allowed")
		}
		positionSet[pos.PositionID()] = struct{}{}
	}

	return nil
}

func NewPlayer(firstName string, lastName string, middleName *string, birthday time.Time, photo *string, cityID uuid.UUID, positions []vo.PlayerPosition, height int64, impactLeg ImpactLeg, marketValue int64) (*Player, error) {
	timeNow := time.Now().UTC()

	player := &Player{
		id:          uuid.New(),
		firstName:   strings.TrimSpace(firstName),
		lastName:    strings.TrimSpace(lastName),
		birthday:    birthday,
		cityID:      cityID,
		positions:   positions,
		height:      height,
		impactLeg:   impactLeg,
		marketValue: marketValue,
		createdAt:   timeNow,
		updatedAt:   timeNow,
	}

	player.SetMiddleName(middleName)
	player.SetPhoto(photo)

	if err := player.validate(); err != nil {
		return nil, err
	}

	return player, nil
}

func NewPlayerWithID(id uuid.UUID, firstName string, lastName string, middleName *string, birthday time.Time, photo *string, cityID uuid.UUID, positions []vo.PlayerPosition, height int64, impactLeg ImpactLeg, marketValue int64, createdAt, updatedAt time.Time) (*Player, error) {
	player := &Player{
		id:          id,
		firstName:   firstName,
		lastName:    lastName,
		birthday:    birthday,
		cityID:      cityID,
		positions:   positions,
		height:      height,
		impactLeg:   impactLeg,
		marketValue: marketValue,
		createdAt:   createdAt,
		updatedAt:   updatedAt,
	}

	player.SetMiddleName(middleName)
	player.SetPhoto(photo)

	if err := player.validate(); err != nil {
		return nil, err
	}

	return player, nil
}
