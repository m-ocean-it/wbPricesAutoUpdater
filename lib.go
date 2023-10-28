package main

import (
	"fmt"
	"log"
	"wbPricesAutoUpdater/domain"
	"wbPricesAutoUpdater/infrastructure"
)

func getTargetPrices() (domain.CatalogPricing, error) {
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
) error {

	var err error

	err = wbClient.UpdatePrices(pricesToSet)
	if err != nil {
		return fmt.Errorf("could not update prices: %w", err)
	}
	log.Println("Updated prices")

	err = wbClient.UpdateDiscounts(discountsToSet)
	if err != nil {
		return fmt.Errorf("could not update discounts: %w", err)
	}
	log.Println("Updated discounts")

	return nil
}
