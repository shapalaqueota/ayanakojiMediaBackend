package bitpay

import (
	"backend/config"
	"bytes"
	"encoding/json"
	"io"
	"net/http"
)

type Invoice struct {
	URL string `json:"url"`
	Id  string `json:"id"`
}

func CreateInvoice(price float64, currency string) (*Invoice, error) {
	url := "https://bitpay.com/invoices"
	body := map[string]interface{}{
		"price":    price,
		"currency": currency,
		"token":    config.AppConfig.BitPayAPIKey,
	}

	jsonBody, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var invoice Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}

type InvoiceStatus struct {
	Status string `json:"status"`
}

func GetInvoiceStatus(invoiceId string) (*InvoiceStatus, error) {
	url := "https://bitpay.com/invoices/" + invoiceId
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+config.AppConfig.BitPayAPIKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			panic(err)
		}
	}(resp.Body)

	var status InvoiceStatus
	if err := json.NewDecoder(resp.Body).Decode(&status); err != nil {
		return nil, err
	}

	return &status, nil
}
