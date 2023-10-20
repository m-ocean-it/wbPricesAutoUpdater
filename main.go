package main

import (
	"context"
	"log"
	"time"
)

func main() {
	var ctx context.Context
	cancel := func() {}

	for {
		currentPrices, err := getCurrentPrices()
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
