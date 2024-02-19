package storage

type IStorage interface {
	UpdateStatus(orderId int)
	GetOutbox(limit int) ([]int, error)
	CreateOutbox(orderId int) int
}
