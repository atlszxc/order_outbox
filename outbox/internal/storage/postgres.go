package storage

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PGStorage struct {
	Conn *pgxpool.Pool
}

func (pgs *PGStorage) UpdateStatus(orderId int) {
	q := `
	UPDATE outbox
	SET is_delivered = true
	WHERE order_id = $1
	`

	pgs.Conn.QueryRow(context.TODO(), q, &orderId)
}

func (pgs *PGStorage) GetOutbox(limit int) ([]int, error) {
	q := `
		SELECT order_id
		FROM outbox
		WHERE is_delivered = false
		LIMIT $1
	`

	rows, err := pgs.Conn.Query(context.TODO(), q, &limit)
	defer rows.Close()
	if err != nil {
		panic(err)
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
func (pgs *PGStorage) CreateOutbox(orderId int) int {
	q := `
		INSERT INTO outbox (order_id)
		VALUES ($1)
		RETURNING id
	`
	var id int
	err := pgs.Conn.QueryRow(context.TODO(), q, &orderId).Scan(&id)
	if err != nil {
		panic(err)
	}

	return id
}

func GetStorage(connStr string) *PGStorage {
	conn, err := pgxpool.New(context.TODO(), connStr)
	if err != nil {
		panic(err)
	}

	return &PGStorage{
		Conn: conn,
	}
}
