package main

import (
	"context"
	"log"
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
) error {
	// TODO: write proper implementation

	log.Printf("executing plan... prices to set: %v, discounts to set: %v\n", pricesToSet, discountsToSet)
	return nil
}
