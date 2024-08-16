package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
)

type Position struct {
	id        uuid.UUID
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func (p *Position) ID() uuid.UUID {
	return p.id
}

func (p *Position) Name() string {
	return p.name
}

func (p *Position) CreatedAt() time.Time {
	return p.createdAt
}

func (p *Position) UpdatedAt() time.Time {
	return p.updatedAt
}

func (p *Position) SetName(name string) error {
	p.name = strings.TrimSpace(name)
	p.updatedAt = time.Now().UTC()

	return p.validate()
}

func (p *Position) validate() error {
	nameLength := len([]rune(p.name))

	if nameLength == 0 {
		return domainerror.Invalid.New("PositionNameIsEmpty", "position name is empty")
	}

	if nameLength < 2 {
		return domainerror.Invalid.New("PositionNameIsTooShort", "position name is too short")
	}

	if p.createdAt.After(p.updatedAt) {
		return domainerror.Invalid.New("InvalidTimeRange", "created_at must be before updated_at")
	}

	return nil
}

func NewPosition(name string) (*Position, error) {
	timeNow := time.Now().UTC()

	position := &Position{
		id:        uuid.New(),
		name:      strings.TrimSpace(name),
		createdAt: timeNow,
		updatedAt: timeNow,
	}

	if err := position.validate(); err != nil {
		return nil, err
	}

	return position, nil
}

func NewPositionWithID(id uuid.UUID, name string, createdAt time.Time, updatedAt time.Time) (*Position, error) {
	position := &Position{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if err := position.validate(); err != nil {
		return nil, err
	}

	return position, nil
}
