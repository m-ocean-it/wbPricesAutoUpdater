package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"wbPricesAutoUpdater/domain"
)

const (
	GET_PRICES_URL    = "https://suppliers-api.wildberries.ru/public/api/v1/info?quantity=0"
	UPDATE_PRICES_URL = "https://suppliers-api.wildberries.ru/public/api/v1/prices"
)

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
	client := &http.Client{}
	req, err := http.NewRequest("GET", GET_PRICES_URL, nil)
	if err != nil {
		return []WbPricingItem{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", c.authorizationHeaderValue())

	response, err := client.Do(req)
	if err != nil {
		return []WbPricingItem{}, fmt.Errorf("error when requesting current prices from Wildberries API: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return []WbPricingItem{}, fmt.Errorf("received non-OK status code (%d) and could not read the body", response.StatusCode)
		}
		return []WbPricingItem{}, fmt.Errorf("response non-OK status code (%d). body: %s", response.StatusCode, string(body))
	}

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

func (c WbOpenApiClient) UpdatePrices(pricesToSet domain.PricesUpdatePlan) error {
	fmt.Printf("prices to set: %v\n", pricesToSet)

	type requestEntry struct {
		NmId  domain.ProductId `json:"nmId"`
		Price domain.Price     `json:"price"`
	}

	requestBody := []requestEntry{}

	for k, v := range pricesToSet {
		entry := requestEntry{
			NmId:  k,
			Price: v,
		}

		requestBody = append(requestBody, entry)
	}

	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("could not marshall the request data into JSON: %w", err)
	}

	req, err := http.NewRequest("POST", UPDATE_PRICES_URL, bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Set("Authorization", c.authorizationHeaderValue())
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return fmt.Errorf("received non-OK status code (%d\n) and could not read the body", response.StatusCode)
		}
		return fmt.Errorf("received non-OK status code (%d\n). body: %s", response.StatusCode, string(body))
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	log.Printf("Updated Wildberries prices. Response body: %s", string(body))

	return nil
}

func (c WbOpenApiClient) UpdateDiscounts(discountsToSet domain.DiscountsUpdatePlan) error {
	return nil
}

func (c WbOpenApiClient) authorizationHeaderValue() string {
	return "Bearer " + c.authToken
}
