package handler

import (
	"awesomeProject/pkg/handler/model"
	"github.com/gin-gonic/gin"
	"net/http"
	"slices"
	"strconv"
)

func (h *Handler) AddTransaction(c *gin.Context) {
	isCorrectSourceType := slices.Contains(sourceType, c.GetHeader("Source-Type"))
	if !isCorrectSourceType {
		newErrorResponse(c, http.StatusBadRequest, "invalid header Source-Type")
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

	transactionWithBalance, err := h.services.AddTransaction(transaction)
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"transaction": model.Transaction{
			State:         transactionWithBalance.State,
			Amount:        strconv.FormatFloat(transactionWithBalance.Amount, 'f', -1, 64),
			TransactionID: transactionWithBalance.TransactionID,
		},
		"user_balance": transactionWithBalance.UserBalance,
	})
}

func (h *Handler) CancelTransactions(c *gin.Context) {

	err := h.services.CancelLatestOddTransactions()
	if err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"transaction": "ok",
	})
}
