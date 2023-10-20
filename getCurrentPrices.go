package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
)

const GET_PRICES_URL = "https://suppliers-api.wildberries.ru/public/api/v1/info?quantity=0"
const authToken = ""

type wbPricingItem struct {
	NmId     uint64   `json:"nmId"`
	Price    price    `json:"price"`
	Discount discount `json:"discount"`
}

func getCurrentPrices() (catalogPricing, error) {
	pricingItems, err := fetchWbPricingItems()
	if err != nil {
		return catalogPricing{}, fmt.Errorf("error when obtaining current Wildberries prices: %w", err)
	}

	pricing := convertWbPricingItemsToCatalogPricing(pricingItems)

	return pricing, nil
}

func fetchWbPricingItems() ([]wbPricingItem, error) {
	// Prepare request
	client := &http.Client{}
	req, err := http.NewRequest("GET", GET_PRICES_URL, nil)
	if err != nil {
		return []wbPricingItem{}, fmt.Errorf("error creating request: %w", err)
	}
	req.Header.Add("Authorization", "Bearer "+authToken)

	// Send request
	response, err := client.Do(req)
	if err != nil {
		return []wbPricingItem{}, fmt.Errorf("error when requesting current prices from Wildberries API: %w", err)
	}
	defer response.Body.Close()

	// Check response's status code
	if response.StatusCode != http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			return []wbPricingItem{}, fmt.Errorf("response has wrong HTTP status (%d)", response.StatusCode)
		}
		return []wbPricingItem{}, fmt.Errorf("response has wrong HTTP status (%d). message: %s", response.StatusCode, string(body))
	}

	// Parse response's JSON
	var pricingItems []wbPricingItem
	body, err := io.ReadAll(response.Body)
	if err != nil {
		return []wbPricingItem{}, errors.New("could not read the response's body")
	}
	err = json.Unmarshal(body, &pricingItems)
	if err != nil {
		return []wbPricingItem{}, errors.New("could not parse the response's JSON")
	}

	return pricingItems, nil
}

func convertWbPricingItemsToCatalogPricing(items []wbPricingItem) catalogPricing {
	catalog := catalogPricing{}
	for _, item := range items {
		catalog[productId(item.NmId)] = pricePair{
			price:    item.Price,
			discount: item.Discount,
		}
	}

	return catalog
}
