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
