package domain

func CompareCurrentVsTargetPrices(current CatalogPricing, target CatalogPricing) (
	PricesUpdatePlan, DiscountsUpdatePlan, error,
) {
	pricesToSet := PricesUpdatePlan{}
	discountsToSet := DiscountsUpdatePlan{}

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
