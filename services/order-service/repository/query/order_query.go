package query

import (
	"context"
	"fmt"
	"log"

	"github.com/daffaromero/retries/services/common/genproto/event"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type OrderQuery interface {
	CreateOrder(context.Context, *event.EventRequest) (*event.EventResponse, error)
	GetOrder(context.Context, *event.GetEventFilter) (*event.GetEventResponse, error)
	GetAllOrders(context.Context, int, int) ([]*event.GetEventResponse, error)
}

type OrderQueryImpl struct {
	Db *pgxpool.Pool
}

func (repo *OrderQueryImpl) CreateOrder(ctx context.Context, er *event.EventRequest) (*event.EventResponse, error) {
	query := `INSERT INTO orders (id, name) VALUES ($1, $2) RETURNING id`
	id := ""
	err := repo.Db.QueryRow(ctx, query, er.Id, er.Name).Scan(&id)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	return &event.EventResponse{Id: er.Id, Success: true}, nil
}

func (repo *OrderQueryImpl) GetOrder(ctx context.Context, ef *event.GetEventFilter) (*event.GetEventResponse, error) {
	var name string
	query := `SELECT name from orders where id=$1`
	err := repo.Db.QueryRow(ctx, query, ef.Id).Scan(&name)
	if err == pgx.ErrNoRows {
		return nil, fmt.Errorf("no records found")
	}
	if err != nil {
		return nil, fmt.Errorf("GetOrder: Bad input :: %e", err)
	} else {
		log.Println(ef.Id, name)
	}
	return &event.GetEventResponse{Id: ef.Id, Name: name}, nil
}

func (repo *OrderQueryImpl) GetAllOrders(ctx context.Context, count, start int) ([]*event.GetEventResponse, error) {
	query := `SELECT name, id from orders LIMIT $1 OFFSET $2`
	rows, err := repo.Db.Query(ctx, query, count, start)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	GetEventResponses := []*event.GetEventResponse{}

	for rows.Next() {
		ge := &event.GetEventResponse{} // Initialize ge variable
		if err := rows.Scan(&ge.Id, &ge.Name); err != nil {
			return nil, err
		}
		GetEventResponses = append(GetEventResponses, ge)
	}
	return GetEventResponses, nil
}
