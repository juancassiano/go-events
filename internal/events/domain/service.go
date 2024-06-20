package domain

import (
	"errors"
	"fmt"
)

var (
	errSpotQuantityMustBePositive = errors.New("Spot quantity must be greather than 0")
)

type spotService struct{}

func NewSpotService() *spotService {
	return &spotService{}
}

func (s *spotService) GenerateSpots(event *Event, quantity int) error {
	if quantity <= 0 {
		return errSpotQuantityMustBePositive
	}

	for i := range quantity {
		spotName := fmt.Sprintf("%c %d", 'A'+1/10, i%10+1)
		spot, err := NewSpot(event, spotName)
		if err != nil {
			return err
		}
		event.Spots = append(event.Spots, *spot)
	}
	return nil
}
