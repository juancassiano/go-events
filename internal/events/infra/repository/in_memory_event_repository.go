package repository

import (
	"errors"
	"sync"

	"github.com/juancassiano/golang/api-sell/internal/events/domain"
)

type inMemoryEventRepository struct {
	mu      sync.RWMutex
	events  map[string]*domain.Event
	spots   map[string]*domain.Spot
	tickets map[string]*domain.Ticket
}

func NewInMemoryEventRepository() domain.EventRepository {
	return &inMemoryEventRepository{
		events:  make(map[string]*domain.Event),
		spots:   make(map[string]*domain.Spot),
		tickets: make(map[string]*domain.Ticket),
	}
}

func (r *inMemoryEventRepository) ListEvents() ([]domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	events := make([]domain.Event, 0, len(r.events))
	for _, event := range r.events {
		events = append(events, *event)
	}
	return events, nil
}

func (r *inMemoryEventRepository) FindEventByID(eventID string) (*domain.Event, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	event, exists := r.events[eventID]
	if !exists {
		return nil, domain.ErrEventNotFound
	}
	return event, nil
}

func (r *inMemoryEventRepository) CreateEvent(event *domain.Event) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.events[event.ID]; exists {
		return errors.New("event already exists")
	}
	r.events[event.ID] = event
	return nil
}

func (r *inMemoryEventRepository) FindSpotByID(spotID string) (*domain.Spot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	spot, exists := r.spots[spotID]
	if !exists {
		return nil, domain.ErrSpotNotFound
	}

	return spot, nil
}

func (r *inMemoryEventRepository) CreateSpot(spot *domain.Spot) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.spots[spot.ID]; exists {
		return errors.New("spot already exists")
	}
	r.spots[spot.ID] = spot

	event, exists := r.events[spot.EventID]
	if !exists {
		return errors.New("event not found")
	}
	event.Spots = append(event.Spots, *spot)
	return nil
}

func (r *inMemoryEventRepository) CreateTicket(ticket *domain.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.tickets[ticket.ID]; exists {
		return errors.New("ticket already exists")
	}
	r.tickets[ticket.ID] = ticket

	event, exists := r.events[ticket.EventID]
	if !exists {
		return errors.New("event not found")
	}
	event.Tickets = append(event.Tickets, *ticket)

	spot, exists := r.spots[ticket.Spot.ID]
	if !exists {
		return errors.New("spot not found")
	}
	spot.TicketID = ticket.ID
	return nil
}

func (r *inMemoryEventRepository) ReserveSpot(spotID, ticketID string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	spot, exists := r.spots[spotID]
	if !exists {
		return domain.ErrSpotNotFound
	}

	spot.Status = domain.SpotStatusSold
	spot.TicketID = ticketID
	return nil
}

func (r *inMemoryEventRepository) FindSpotsByEventID(eventID string) ([]*domain.Spot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	spots := make([]*domain.Spot, 0)
	for _, spot := range r.spots {
		if spot.EventID == eventID {
			spots = append(spots, spot)
		}
	}
	return spots, nil
}

func (r *inMemoryEventRepository) FindSpotByName(eventID, name string) (*domain.Spot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, spot := range r.spots {
		if spot.EventID == eventID && spot.Name == name {
			return spot, nil
		}
	}
	return nil, domain.ErrSpotNotFound
}
