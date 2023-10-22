package main

import (
	"context"
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

	wbClient := infrastructure.NewWbOpenApiClient(wbAuthToken)

	cancelSaveCurrentPricesFunc := func() {}

	for {
		cancelSaveCurrentPricesFunc = run_cycle(wbClient, cancelSaveCurrentPricesFunc)
		sleep()
	}
}

func run_cycle(
	wbClient infrastructure.WbOpenApiClient,
	cancelSaveCurrentPricesFunc context.CancelFunc,
) context.CancelFunc {

	// TODO: cache current prices to avoid constantly fetching them
	currentPrices, err := getCurrentPrices(wbClient)
	if err != nil {
		log.Println(err)
		return nil
	}

	log.Printf("received current prices: %d\n", len(currentPrices))

	// cancel context so that the previous instance of save_current_prices() stops execution
	if cancelSaveCurrentPricesFunc != nil {
		cancelSaveCurrentPricesFunc()
	}
	// new context and cancellation func for new invocation of save_current_prices()
	saveCurrentPricesCtx, cancelSaveCurrentPricesFunc := context.WithCancel(context.TODO())

	go saveCurrentPrices(saveCurrentPricesCtx, currentPrices)

	targetPrices, err := getTargetPrices()
	if err != nil {
		log.Println(err)
		return cancelSaveCurrentPricesFunc
	}

	pricesToSet, discountsToSet, err := compareCurrentVsTargetPrices(currentPrices, targetPrices)
	if err != nil {
		log.Println(err)
		return cancelSaveCurrentPricesFunc
	}

	errs := executePricingUpdatePlan(currentPrices, pricesToSet, discountsToSet, wbClient)
	if len(errs) > 0 {
		for _, e := range errs {
			log.Println(e)
		}
		return cancelSaveCurrentPricesFunc
	}

	return cancelSaveCurrentPricesFunc
}

func sleep() {
	log.Println("sleeping for 10 seconds")
	time.Sleep(time.Second * 10)
	log.Println("===========END=OF=LOOP=============")
}
