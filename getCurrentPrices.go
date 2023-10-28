package main

import (
	"fmt"
	"log"
	"time"
	"wbPricesAutoUpdater/domain"
	"wbPricesAutoUpdater/infrastructure"
)

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

type CurrentPricingCache interface {
	GetPricing() (pricing []infrastructure.WbPricingItem, cacheAge time.Duration, err error)
	SavePricing([]infrastructure.WbPricingItem) error
}

func NewPricingServer(cache CurrentPricingCache, wbClient infrastructure.WbOpenApiClient) PricingServer {
	return PricingServer{cache, wbClient}
}

type PricingServer struct {
	cache    CurrentPricingCache
	wbClient infrastructure.WbOpenApiClient
}

func (ps *PricingServer) getCurrentPrices() ([]infrastructure.WbPricingItem, error) {
	pricingItems, err := ps.wbClient.FetchWbPricingItems()
	if err != nil {
		return []infrastructure.WbPricingItem{}, fmt.Errorf("error when obtaining current Wildberries prices: %w", err)
	}

	return pricingItems, nil
}

type response struct {
	Pricing  domain.CatalogPricing
	CacheAge time.Duration
}

func (ps *PricingServer) FetchAndCacheCurrentPrices() (response, error) {
	log.Println("Fetching current prices")

	currentPricing, err := ps.getCurrentPrices()
	if err != nil {
		log.Println("Could not fetch current prices -> fetching from cache")

		pricing, cacheAge, err := ps.cache.GetPricing()

		return response{
			Pricing:  convertWbPricingItemsToCatalogPricing(pricing),
			CacheAge: cacheAge,
		}, err
	}

	log.Println("Got current prices")

	go func() {
		log.Println("Saving current prices")

		ps.cache.SavePricing(currentPricing)

		log.Println("Saved current prices")
	}()

	return response{
		Pricing:  convertWbPricingItemsToCatalogPricing(currentPricing),
		CacheAge: 0,
	}, nil
}
