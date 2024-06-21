package repository

import (
	"database/sql"

	"github.com/juancassiano/golang/api-sell/internal/events/domain"
)

type mysqlEventRepository struct {
	db *sql.DB
}

// func NewMysqlEventRepository(db *sql.DB) (*domain.EventRepository, error) {
// 	return &mysqlEventRepository{db}, nil
// }

func (r *mysqlEventRepository) CreateSpot(spot *domain.Spot) error {

	query := `
		INSERT INTO spots (id, event_id, name, status, ticket_id)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, spot.ID, spot.EventID, spot.Name, spot.Status, spot.TicketID)
	return err
}

func (r *mysqlEventRepository) ReserveSpot(spotID, ticketID string) error {
	query := `
		UPDATE spots
		SET status = ?, ticket_id = ?
		WHERE id = ?
	`
	_, err := r.db.Exec(query, domain.SpotStatusSold, ticketID, spotID)
	return err
}

func (r *mysqlEventRepository) CreateTicket(ticket domain.Ticket) error {
	query := `
		INSERT INTO tickets (id, event_id, spot_id, ticket_type, price)
		VALUES (?, ?, ?, ?, ?)
	`
	_, err := r.db.Exec(query, ticket.ID, ticket.EventID, ticket.Spot.ID, ticket.TicketType, ticket.Price)
	return err
}
