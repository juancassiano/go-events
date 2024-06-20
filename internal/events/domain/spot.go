package domain

import (
	"errors"

	"github.com/google/uuid"
)

type SpotStatus string

var (
	ErrInvalidSpotNumber             = errors.New("invalid spot number")
	ErrSpotNotFound                  = errors.New("spot not found")
	ErrSpotAlreadyReserved           = errors.New("spot already reserved")
	ErrSpotNameIsRequired            = errors.New("spot name is required")
	ErrSpotNameMustBeGreatherThanTwo = errors.New("spot name must be at least 2 characters long")
	ErrSpotNameMustStartWithLetter   = errors.New("spot name must start with a letter")
	ErrSpotNameMustEndWithNumber     = errors.New("spot name must end with a number")
)

const (
	SpotStatusAvailable SpotStatus = "available"
	SpotStatusSold      SpotStatus = "sold"
)

type Spot struct {
	ID       string
	EventID  string
	Name     string
	Status   SpotStatus
	TicketID string
}

func NewSpot(event *Event, name string) (*Spot, error) {
	spot := &Spot{
		ID:      uuid.New().String(),
		EventID: event.ID,
		Name:    name,
		Status:  SpotStatusAvailable,
	}

	if err := spot.Validate(); err != nil {
		return nil, err
	}
	return spot, nil

}

func (s *Spot) Validate() error {
	if len(s.Name) == 0 {
		return ErrSpotNameIsRequired
	}
	if len(s.Name) < 2 {
		return ErrSpotNameMustBeGreatherThanTwo
	}
	// Validate if the spot name is in the correct format
	if s.Name[0] < 'A' || s.Name[0] > 'Z' {
		return ErrSpotNameMustStartWithLetter
	}
	if s.Name[1] < '0' || s.Name[1] > '9' {
		return ErrSpotNameMustEndWithNumber
	}
	return nil
}
