package main

import "context"

type productId uint64
type price uint16

type productPrices map[productId]price
type pricesUpdatePlan productPrices

func getCurrentPrices() (productPrices, error) {
	return productPrices{}, nil
}

func saveCurrentPrices(ctx context.Context, prices productPrices) error {
	return nil
}

func getTargetPrices() (productPrices, error) {
	return productPrices{}, nil
}

func compareCurrentVsTargetPrices(current productPrices, target productPrices) (pricesUpdatePlan, error) {
	return pricesUpdatePlan{}, nil
}

func executePriceUpdatePlan(tasks pricesUpdatePlan) error {
	return nil
}
