package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
	"wbPricesAutoUpdater/domain"
)

const (
	GET_PRICES_URL       = "https://suppliers-api.wildberries.ru/public/api/v1/info?quantity=0"
	UPDATE_PRICES_URL    = "https://suppliers-api.wildberries.ru/public/api/v1/prices"
	UPDATE_DISCOUNTS_URL = "https://suppliers-api.wildberries.ru/public/api/v1/updateDiscounts"
)

type WbPricingItem struct {
	NmId     uint64          `json:"nmId"`
	Price    domain.Price    `json:"price"`
	Discount domain.Discount `json:"discount"`
}

type WbOpenApiClient struct {
	authToken  string
	httpClient *http.Client
}

func NewWbOpenApiClient(authToken string) WbOpenApiClient {
	return WbOpenApiClient{
		authToken: authToken,
		httpClient: &http.Client{
			Timeout: 15 * time.Second,
		},
	}
}

func (c WbOpenApiClient) FetchWbPricingItems() ([]WbPricingItem, error) {
	req, err := http.NewRequest("GET", GET_PRICES_URL, nil)
	if err != nil {
		return []WbPricingItem{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", c.authorizationHeaderValue())

	response, err := c.httpClient.Do(req)
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
	if len(pricesToSet) == 0 {
		log.Println("No prices to set, therefore, UpdatePrices completes successfully.")
		return nil
	}
	fmt.Printf("Setting Wildberries prices: %v\n", pricesToSet)

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

	req, err := c.newAuthorizedJsonRequest(UPDATE_PRICES_URL, jsonRequestBody)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nonOkErrorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	log.Printf("Updated Wildberries prices. Response body: %s", string(body))

	return nil
}

func (c WbOpenApiClient) UpdateDiscounts(discountsToSet domain.DiscountsUpdatePlan) error {
	if len(discountsToSet) == 0 {
		log.Println("No discounts to set, therefore, UpdateDiscounts completes successfully.")
		return nil
	}
	fmt.Printf("Setting Wildberries discounts: %v\n", discountsToSet)

	type requestEntry struct {
		Nm       domain.ProductId `json:"nm"`
		Discount domain.Discount  `json:"discount"`
	}

	requestBody := []requestEntry{}

	for k, v := range discountsToSet {
		entry := requestEntry{
			Nm:       k,
			Discount: v,
		}

		requestBody = append(requestBody, entry)
	}

	jsonRequestBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("could not marshall the request data into JSON: %w", err)
	}

	req, err := c.newAuthorizedJsonRequest(UPDATE_DISCOUNTS_URL, jsonRequestBody)
	if err != nil {
		return err
	}

	response, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("error making request: %w", err)
	}
	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nonOkErrorFromResponse(response)
	}

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	log.Printf("Updated Wildberries discounts. Response body: %s", string(body))

	return nil
}

func (c WbOpenApiClient) authorizationHeaderValue() string {
	return "Bearer " + c.authToken
}

func (c WbOpenApiClient) newAuthorizedJsonRequest(url string, jsonBody []byte) (
	*http.Request, error,
) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("error creating authorized JSON request: %w", err)
	}
	req.Header.Set("Authorization", c.authToken)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}
