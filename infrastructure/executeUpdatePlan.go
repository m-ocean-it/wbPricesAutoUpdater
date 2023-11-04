package infrastructure

import (
	"fmt"
	"log"
	"wbPricesAutoUpdater/domain"
)

func ExecutePricingUpdatePlan(
	currentCatalogPricing domain.CatalogPricing,
	pricesToSet domain.PricesUpdatePlan,
	discountsToSet domain.DiscountsUpdatePlan,
	wbClient WbOpenApiClient,
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
