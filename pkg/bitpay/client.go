package bitpay

import (
	"backend/config"
	"bytes"
	"encoding/json"
	"net/http"
)

type Invoice struct {
	URL string `json:"url"`
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
	defer resp.Body.Close()

	var invoice Invoice
	if err := json.NewDecoder(resp.Body).Decode(&invoice); err != nil {
		return nil, err
	}

	return &invoice, nil
}
