package main

import (
	"context"
	"log"
	"time"
)

type productId uint64
type price uint16
type discount uint8

type pricePair struct {
	price    price
	discount discount
}
type catalogPricing map[productId]pricePair

type pricesUpdatePlan map[productId]price
type discountsUpdatePlan map[productId]discount

func saveCurrentPrices(ctx context.Context, prices catalogPricing) error {
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

func getTargetPrices() (catalogPricing, error) {
	// return catalogPricing{}, nil
	return catalogPricing{
		1: {price: 380, discount: 10},
		2: {price: 400, discount: 12},
	}, nil // TODO: implement
}

func compareCurrentVsTargetPrices(current catalogPricing, target catalogPricing) (pricesUpdatePlan, discountsUpdatePlan, error) {
	pricesToSet := pricesUpdatePlan{}
	discountsToSet := discountsUpdatePlan{}

	for productId, targetPricePair := range target {
		currentPricePair, ok := current[productId]
		if !ok {
			continue
		}

		if currentPricePair.price != targetPricePair.price {
			pricesToSet[productId] = targetPricePair.price
		}
		if currentPricePair.discount != targetPricePair.discount {
			discountsToSet[productId] = targetPricePair.discount
		}
	}

	return pricesToSet, discountsToSet, nil
}

func executePricingUpdatePlan(
	currentCatalogPricing catalogPricing,
	pricesToSet pricesUpdatePlan,
	discountsToSet discountsUpdatePlan,
) error {
	// TODO: write proper implementation

	log.Printf("executing plan... prices to set: %v, discounts to set: %v\n", pricesToSet, discountsToSet)
	return nil
}
