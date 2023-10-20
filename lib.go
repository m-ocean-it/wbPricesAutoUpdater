package main

import (
	"context"
	"log"
	"sync"
	"time"
	"wbPricesAutoUpdater/domain"
)

type pricesUpdatePlan map[domain.ProductId]domain.Price
type discountsUpdatePlan map[domain.ProductId]domain.Discount

func saveCurrentPrices(ctx context.Context, prices domain.CatalogPricing) error {
	// TODO: write proper implementation

	// log.Printf("saving current prices: %v\n", prices)
	time.Sleep(time.Second * 10)

	if err := ctx.Err(); err != nil {
		log.Println("saving current prices was canceled")
		return err
	}

	log.Println("saved current prices")
	return nil
}

func getTargetPrices() (domain.CatalogPricing, error) {
	// return catalogPricing{}, nil
	return domain.CatalogPricing{
		94640599: {Price: 1665, Discount: 45},
	}, nil // TODO: implement
}

func compareCurrentVsTargetPrices(current domain.CatalogPricing, target domain.CatalogPricing) (pricesUpdatePlan, discountsUpdatePlan, error) {
	pricesToSet := pricesUpdatePlan{}
	discountsToSet := discountsUpdatePlan{}

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
	pricesToSet pricesUpdatePlan,
	discountsToSet discountsUpdatePlan,
) []error {
	var mu sync.Mutex
	encounteredErrors := []error{}

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		err := updatePrices(pricesToSet)
		if err != nil {
			mu.Lock()
			encounteredErrors = append(encounteredErrors, err)
			mu.Unlock()
		}
		wg.Done()
	}()
	go func() {
		err := updateDiscounts(discountsToSet)
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

func updatePrices(pricesToSet pricesUpdatePlan) error          { return nil }
func updateDiscounts(discountsToSet discountsUpdatePlan) error { return nil }
