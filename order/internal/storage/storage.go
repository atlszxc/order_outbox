package storage

import (
	"order/internal/order/dtos"
	"order/internal/order/entity"
)

type Storage interface {
	GetCompletedOrderId(limit int) ([]int, error)
	Create(dto dtos.CreateOrderDto) (*entity.Order, error)
	UpdateDeliveredStatus(id int, status bool)
	UpdateStatus(dto dtos.UpdateStatus) error
}
