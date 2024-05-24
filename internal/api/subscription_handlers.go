package api

import (
	"backend/pkg/bitpay"
	"github.com/gin-gonic/gin"
	"net/http"
)

type CreateInvoiceRequest struct {
	Price    float64 `json:"price"`
	Currency string  `json:"currency"`
}

func CreateInvoice(c *gin.Context) {
	var request CreateInvoiceRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	invoice, err := bitpay.CreateInvoice(request.Price, request.Currency)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, invoice)
}
