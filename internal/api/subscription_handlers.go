package api

import (
	"backend/pkg/bitpay"
	"github.com/gin-gonic/gin"
	"net/http"
)

func CreateInvoice(c *gin.Context) {
	invoice, err := bitpay.CreateInvoice(9.99, "USD")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"invoiceUrl": invoice.URL, "invoiceId": invoice.Id})
}

func GetInvoiceStatus(c *gin.Context) {
	invoiceId := c.Param("id")
	status, err := bitpay.GetInvoiceStatus(invoiceId)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": status.Status})
}
