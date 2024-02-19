package dtos

type CreateOrderDto struct {
	Products []string
}

type UpdateStatus struct {
	Id        int
	Confirmed bool
}

type ConfirmedOrderId struct {
	Id int
}
