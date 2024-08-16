package domain

import (
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/javascriptizer1/tm-shared.backend/domainerror"
)

type City struct {
	id        uuid.UUID
	name      string
	createdAt time.Time
	updatedAt time.Time
}

func (c *City) ID() uuid.UUID {
	return c.id
}

func (c *City) Name() string {
	return c.name
}

func (c *City) CreatedAt() time.Time {
	return c.createdAt
}

func (c *City) UpdatedAt() time.Time {
	return c.updatedAt
}

func (c *City) SetName(name string) error {
	c.name = strings.TrimSpace(name)
	c.updatedAt = time.Now().UTC()

	return c.validate()
}

func (c *City) validate() error {
	nameLength := len([]rune(c.name))

	if nameLength == 0 {
		return domainerror.Invalid.New("CityNameIsEmpty", "city name is empty")
	}

	if nameLength < 2 {
		return domainerror.Invalid.New("CityNameIsTooShort", "city name is too short")
	}

	if c.createdAt.After(c.updatedAt) {
		return domainerror.Invalid.New("InvalidTimeRange", "created_at must be before updated_at")
	}

	return nil
}

func NewCity(name string) (*City, error) {
	timeNow := time.Now().UTC()

	city := &City{
		id:        uuid.New(),
		name:      strings.TrimSpace(name),
		createdAt: timeNow,
		updatedAt: timeNow,
	}

	if err := city.validate(); err != nil {
		return nil, err
	}

	return city, nil
}

func NewCityWithID(id uuid.UUID, name string, createdAt time.Time, updatedAt time.Time) (*City, error) {
	city := &City{
		id:        id,
		name:      name,
		createdAt: createdAt,
		updatedAt: updatedAt,
	}

	if err := city.validate(); err != nil {
		return nil, err
	}

	return city, nil
}
