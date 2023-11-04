package main

import (
	"fmt"
	"log"
	"os"
	"time"
	"wbPricesAutoUpdater/domain"
	"wbPricesAutoUpdater/infrastructure"
)

const WB_OPENAPI_AUTH_TOKEN_ENV_VAR = "WB_OPENAPI_AUTH_TOKEN"
const MAX_PRICES_CACHE_AGE = 15 * time.Minute

func main() {
	wbAuthToken := os.Getenv(WB_OPENAPI_AUTH_TOKEN_ENV_VAR)
	if wbAuthToken == "" {
		log.Fatalf("%s environment variable must be set\n", WB_OPENAPI_AUTH_TOKEN_ENV_VAR)
	}

	cache := infrastructure.JsonCurrentPricingCache{
		FileHandler: infrastructure.NewConcurrentFileHandler(
			"./pricing_cache.json",
		),
	}

	wbClient := infrastructure.NewWbOpenApiClient(wbAuthToken)

	pricingServer := infrastructure.NewPricingServer(&cache, wbClient)

	var err error
	for {
		err = run_cycle(wbClient, pricingServer)
		if err != nil {
			log.Printf("Cycle did not complete successfully. Error: %s\n", err)
		} else {
			log.Println("Cycle completed successfully")
		}

		sleep()
	}
}

func run_cycle(wbClient infrastructure.WbOpenApiClient, pricingServer infrastructure.PricingServer) error {
	resp, err := pricingServer.FetchAndCacheCurrentPrices()
	if err != nil {
		return err
	}

	currentPrices := resp.Pricing
	currentPricesCacheAge := resp.CacheAge

	if currentPricesCacheAge > MAX_PRICES_CACHE_AGE {
		return fmt.Errorf("prices cache too old: %s. Max value: %s",
			currentPricesCacheAge,
			MAX_PRICES_CACHE_AGE)
	}

	log.Printf("Received prices: %d. Cache age: %s\n", len(currentPrices), resp.CacheAge)

	targetPrices, err := infrastructure.GetTargetPrices()
	if err != nil {
		log.Println(err)
		return err
	}

	pricesToSet, discountsToSet, err := domain.CompareCurrentVsTargetPrices(currentPrices, targetPrices)
	if err != nil {
		log.Println(err)
		return err
	}

	err = infrastructure.ExecutePricingUpdatePlan(currentPrices, pricesToSet, discountsToSet, wbClient)

	return err
}

func sleep() {
	log.Println("sleeping for 10 seconds")
	time.Sleep(time.Second * 10)
	log.Println("===========END=OF=LOOP=============")
}
