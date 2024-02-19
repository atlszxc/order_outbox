package entity

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Order struct {
	Id          int
	Products    []string
	Date        pgtype.Date
	Confirmed   bool
	IsDelivered bool
}
