package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"wbPricesAutoUpdater/domain"
)

const GET_PRICES_URL = "https://suppliers-api.wildberries.ru/public/api/v1/info?quantity=0"

type WbPricingItem struct {
	NmId     uint64          `json:"nmId"`
	Price    domain.Price    `json:"price"`
	Discount domain.Discount `json:"discount"`
}

type WbOpenApiClient struct {
	authToken string
}

func NewWbOpenApiClient(token string) WbOpenApiClient {
	return WbOpenApiClient{token}
}

func (c WbOpenApiClient) FetchWbPricingItems() ([]WbPricingItem, error) {
	// Prepare request
	client := &http.Client{}
	req, err := http.NewRequest("GET", GET_PRICES_URL, nil)
	if err != nil {
		return []WbPricingItem{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+c.authToken)

	// Send request
	response, err := client.Do(req)
	if err != nil {
		return []WbPricingItem{}, fmt.Errorf("error when requesting current prices from Wildberries API: %w", err)
	}
	defer response.Body.Close()

	// Check response's status code
	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return []WbPricingItem{}, fmt.Errorf("response has wrong HTTP status (%d)", response.StatusCode)
		}
		return []WbPricingItem{}, fmt.Errorf("response has wrong HTTP status (%d). message: %s", response.StatusCode, string(body))
	}

	// Parse response's JSON
	var pricingItems []WbPricingItem
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []WbPricingItem{}, errors.New("could not read the response's body")
	}
	err = json.Unmarshal(body, &pricingItems)
	if err != nil {
		return []WbPricingItem{}, errors.New("could not parse the response's JSON")
	}

	return pricingItems, nil
}
