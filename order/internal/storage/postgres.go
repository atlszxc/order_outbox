package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgtype"
	"github.com/jackc/pgx/v5/pgxpool"
	"order/internal/order/dtos"
	"order/internal/order/entity"
)

type PGStorage struct {
	Conn *pgxpool.Pool
}

func (pgs *PGStorage) GetCompletedOrderId(limit int) ([]int, error) {
	q := `
		SELECT id
		FROM "order"
		WHERE confirmed = true AND is_delivered = false
		LIMIT $1
	`

	rows, err := pgs.Conn.Query(context.TODO(), q, &limit)
	if err != nil {
		return nil, err
	}

	var ids []int

	for rows.Next() {
		var id int
		err := rows.Scan(&id)
		if err != nil {
			return ids, err
		}

		ids = append(ids, id)
	}

	return ids, nil
}

func (pgs *PGStorage) Create(dto dtos.CreateOrderDto) (*entity.Order, error) {
	q := `
		INSERT INTO "order" (products)
		VALUES ($1)
		RETURNING id, products, date, confirmed
	`

	defer pgs.Conn.Close()

	var res entity.Order
	var dt pgtype.Date

	err := pgs.Conn.QueryRow(context.TODO(), q, &dto.Products).Scan(&res.Id, &res.Products, &dt, &res.Confirmed)
	if err != nil {
		return nil, err
	}
	res.Date = dt
	return &res, nil
}

func (pgs *PGStorage) UpdateDeliveredStatus(id int, status bool) {
	q := `
		UPDATE "order"
		SET is_delivered = $2
		WHERE id = $1
	`

	pgs.Conn.QueryRow(context.TODO(), q, &id, &status)
}

func (pgs *PGStorage) UpdateStatus(dto dtos.UpdateStatus) error {
	defer pgs.Conn.Close()
	q := `
		UPDATE "order"
		SET confirmed = $2
		WHERE id = $1
	`

	rows, err := pgs.Conn.Query(context.TODO(), q, &dto.Id, &dto.Confirmed)
	defer rows.Close()
	if err != nil {
		return err
	}

	return nil
}

func InitStorage(connStr string) (*PGStorage, error) {
	conn, err := pgxpool.New(context.TODO(), connStr)
	if err != nil {
		return nil, err
	}

	return &PGStorage{
		Conn: conn,
	}, nil
}
