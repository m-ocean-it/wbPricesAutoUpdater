package main

import (
	"log"
	"os"
	"time"
	"wbPricesAutoUpdater/infrastructure"
)

const WB_OPENAPI_AUTH_TOKEN_ENV_VAR = "WB_OPENAPI_AUTH_TOKEN"

func main() {
	wbAuthToken := os.Getenv(WB_OPENAPI_AUTH_TOKEN_ENV_VAR)
	if wbAuthToken == "" {
		log.Fatalf("%s environment variable must be set\n", WB_OPENAPI_AUTH_TOKEN_ENV_VAR)
	}

	cache := infrastructure.NewJsonCurrentPricingCache("./pricing_cache.json")
	wbClient := infrastructure.NewWbOpenApiClient(wbAuthToken)

	pricingServer := NewPricingServer(&cache, wbClient)

	var ok bool
	for {
		ok = run_cycle(wbClient, pricingServer)
		if ok {
			log.Println("Cycle completed successfully")
		} else {
			log.Println("Cycle did not complete successfully")
		}
		sleep()
	}
}

func run_cycle(wbClient infrastructure.WbOpenApiClient, pricingServer PricingServer) bool {

	resp, err := pricingServer.FetchAndCacheCurrentPrices()
	if err != nil {
		log.Println(err)
		return false
	}

	currentPrices := resp.Pricing

	log.Printf("Received prices: %d. Cache age: %s\n", len(currentPrices), resp.CacheAge)

	targetPrices, err := getTargetPrices()
	if err != nil {
		log.Println(err)
		return false
	}

	pricesToSet, discountsToSet, err := compareCurrentVsTargetPrices(currentPrices, targetPrices)
	if err != nil {
		log.Println(err)
		return false
	}

	errs := executePricingUpdatePlan(currentPrices, pricesToSet, discountsToSet, wbClient)
	if len(errs) > 0 {
		for _, e := range errs {
			log.Println(e)
		}
		return false
	}

	return true
}

func sleep() {
	log.Println("sleeping for 10 seconds")
	time.Sleep(time.Second * 10)
	log.Println("===========END=OF=LOOP=============")
}
