package main

import (
	"context"
	"log"
	"sync"
	"time"
	"wbPricesAutoUpdater/domain"
	"wbPricesAutoUpdater/infrastructure"
)

func saveCurrentPrices(ctx context.Context, prices domain.CatalogPricing) error {
	// TODO: write proper implementation

	log.Println("saving current prices")
	time.Sleep(time.Second * 20)

	if err := ctx.Err(); err != nil {
		log.Println("saving current prices was canceled")
		return err
	}

	log.Println("saved current prices")
	return nil
}

func getTargetPrices() (domain.CatalogPricing, error) {
	// return catalogPricing{}, nil
	return domain.CatalogPricing{}, nil // TODO: implement
}

func compareCurrentVsTargetPrices(current domain.CatalogPricing, target domain.CatalogPricing) (domain.PricesUpdatePlan, domain.DiscountsUpdatePlan, error) {
	pricesToSet := domain.PricesUpdatePlan{}
	discountsToSet := domain.DiscountsUpdatePlan{}

	for productId, targetPricePair := range target {
		currentPricePair, ok := current[productId]
		if !ok {
			continue
		}

		if currentPricePair.Price != targetPricePair.Price {
			pricesToSet[productId] = targetPricePair.Price
		}
		if currentPricePair.Discount != targetPricePair.Discount {
			discountsToSet[productId] = targetPricePair.Discount
		}
	}

	return pricesToSet, discountsToSet, nil
}

func executePricingUpdatePlan(
	currentCatalogPricing domain.CatalogPricing,
	pricesToSet domain.PricesUpdatePlan,
	discountsToSet domain.DiscountsUpdatePlan,
	wbClient infrastructure.WbOpenApiClient,
) []error {
	var mu sync.Mutex
	encounteredErrors := []error{}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		err := wbClient.UpdatePrices(pricesToSet)
		if err != nil {
			mu.Lock()
			encounteredErrors = append(encounteredErrors, err)
			mu.Unlock()
		}
		wg.Done()
	}()
	go func() {
		err := wbClient.UpdateDiscounts(discountsToSet)
		if err != nil {
			mu.Lock()
			encounteredErrors = append(encounteredErrors, err)
			mu.Unlock()
		}
		wg.Done()
	}()
	wg.Wait()

	return nil
}
