package vo

import "github.com/google/uuid"

type PlayerPosition struct {
	positionID uuid.UUID
	main       bool
}

func (pp *PlayerPosition) PositionID() uuid.UUID {
	return pp.positionID
}

func (pp *PlayerPosition) Main() bool {
	return pp.main
}

func NewPlayerPosition(positionID uuid.UUID, main bool) PlayerPosition {
	return PlayerPosition{
		positionID: positionID,
		main:       main,
	}
}
