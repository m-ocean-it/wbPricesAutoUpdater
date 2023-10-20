package main

import (
	"context"
	"log"
	"os"
	"time"
)

const WB_OPENAPI_AUTH_TOKEN_LABEL = "WB_OPENAPI_AUTH_TOKEN"

func main() {
	wbAuthToken := os.Getenv(WB_OPENAPI_AUTH_TOKEN_LABEL)
	if wbAuthToken == "" {
		log.Fatalf("%s environment variable must be set\n", WB_OPENAPI_AUTH_TOKEN_LABEL)
	}

	wbClient := NewWbOpenApiClient(wbAuthToken)

	var ctx context.Context
	cancel := func() {}

	for {
		currentPrices, err := getCurrentPrices(wbClient)
		if err != nil {
			log.Println(err)
			sleep()
			continue
		}

		// cancel context so that the previous instance of save_current_prices() stops execution
		cancel()
		// new context and cancellation func for new invocation of save_current_prices()
		ctx, cancel = context.WithCancel(context.TODO())

		go saveCurrentPrices(ctx, currentPrices)

		targetPrices, err := getTargetPrices()
		if err != nil {
			log.Println(err)
			sleep()
			continue
		}

		pricesToSet, discountsToSet, err := compareCurrentVsTargetPrices(currentPrices, targetPrices)
		if err != nil {
			log.Println(err)
			sleep()
			continue
		}

		err = executePricingUpdatePlan(currentPrices, pricesToSet, discountsToSet)
		if err != nil {
			log.Println(err)
			sleep()
			continue
		}

		sleep()
	}
}

func sleep() {
	log.Println("sleeping for 10 seconds")
	time.Sleep(time.Second * 4)
	log.Println("===========END=OF=LOOP=============")
}
