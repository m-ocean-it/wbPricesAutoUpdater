package main

import (
	"fmt"
	"wbPricesAutoUpdater/domain"
	"wbPricesAutoUpdater/infrastructure"
)

func getCurrentPrices(wbClient infrastructure.WbOpenApiClient) (domain.CatalogPricing, error) {
	pricingItems, err := wbClient.FetchWbPricingItems()
	if err != nil {
		return domain.CatalogPricing{}, fmt.Errorf("error when obtaining current Wildberries prices: %w", err)
	}

	pricing := convertWbPricingItemsToCatalogPricing(pricingItems)

	return pricing, nil
}

func convertWbPricingItemsToCatalogPricing(items []infrastructure.WbPricingItem) domain.CatalogPricing {
	catalog := domain.CatalogPricing{}
	for _, item := range items {
		catalog[domain.ProductId(item.NmId)] = domain.PricePair{
			Price:    item.Price,
			Discount: item.Discount,
		}
	}

	return catalog
}
