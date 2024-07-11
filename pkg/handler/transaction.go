package handler

import (
	"net/http"
	"slices"
	"strconv"

	"userAccountBalanceService/pkg/handler/model"
	servicemodel "userAccountBalanceService/pkg/service/model"

	"github.com/gin-gonic/gin"
)

func (h *Handler) AddTransaction(c *gin.Context) {
	isCorrectSourceType := slices.Contains(sourceType, c.GetHeader("Source-Type"))
	if !isCorrectSourceType {
		newErrorResponse(c, http.StatusBadRequest, "invalid Source-Type header")
		return
	}
	var transaction model.Transaction
	if err := c.BindJSON(&transaction); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	isCorrectState := slices.Contains(states, transaction.State)
	if !isCorrectState {
		newErrorResponse(c, http.StatusBadRequest, "invalid state")
		return
	}

	transactionWithBalance, err := h.service.AddTransaction(transaction)
	if err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	c.JSON(http.StatusOK, createResponse(transactionWithBalance))
}

func createResponse(transactionWithBalance *servicemodel.TransactionWithBalance) map[string]interface{} {
	return map[string]interface{}{
		"transaction": model.Transaction{
			State:         transactionWithBalance.State,
			Amount:        strconv.FormatFloat(transactionWithBalance.Amount, 'f', -1, 64),
			TransactionID: transactionWithBalance.TransactionID,
		},
		"user_balance": transactionWithBalance.UserBalance,
	}
}
