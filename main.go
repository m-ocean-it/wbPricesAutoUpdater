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

	var ctx context.Context
	cancel := func() {}

	for {
		// TODO: cache current prices to avoid constantly fetching them
		currentPrices, err := getCurrentPrices(wbClient)
		if err != nil {
			log.Println(err)
			sleep()
			continue
		}

		log.Printf("received current prices: %d\n", len(currentPrices))

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

		errs := executePricingUpdatePlan(currentPrices, pricesToSet, discountsToSet, wbClient)
		if len(errs) > 0 {
			for _, e := range errs {
				log.Println(e)
			}
			sleep()
			continue
		}

		sleep()
	}
}

func sleep() {
	log.Println("sleeping for 10 seconds")
	time.Sleep(time.Second * 10)
	log.Println("===========END=OF=LOOP=============")
}
