package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"order/internal/order/dtos"
	"order/internal/storage"
)

func CreateOrderHandler(ctx *gin.Context) {
	var dto dtos.CreateOrderDto

	serializeErr := ctx.BindJSON(&dto)
	if serializeErr != nil {
		panic(serializeErr)
	}

	s, err := storage.InitStorage("postgres://root:postgres@orderdb:5432/order")
	if err != nil {
		panic(err)
	}

	order, dbErr := s.Create(dto)
	if dbErr != nil {
		panic(dbErr)
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"order": order,
	})
}

func UpdateStatusHandler(ctx *gin.Context) {
	var dto dtos.UpdateStatus
	serializeErr := ctx.BindJSON(&dto)
	if serializeErr != nil {
		panic(serializeErr)
	}
	s, err := storage.InitStorage("postgres://root:postgres@orderdb:5432/order")
	if err != nil {
		panic(err)
	}

	dbErr := s.UpdateStatus(dto)
	if dbErr != nil {
		panic(dbErr)
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})

}
