package handler

import (
	"github.com/gin-gonic/gin"

	"userAccountBalanceService/pkg/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	api := router.Group("/api")
	{
		transaction := api.Group("/transaction")
		{
			transaction.POST("/", h.AddTransaction)
		}
	}

	return router
}
