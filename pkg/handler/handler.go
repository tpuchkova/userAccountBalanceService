package handler

import (
	"awesomeProject/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		transaction := api.Group("/transaction")
		{
			transaction.POST("/", h.AddTransaction)
			transaction.PUT("/", h.CancelTransactions)
		}
	}

	return router
}
