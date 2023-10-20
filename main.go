package main

import "context"

func main() {
	var ctx context.Context
	var cancel context.CancelFunc

	for {
		currentPrices, err := getCurrentPrices()
		if err != nil {
			// TODO: log err
			continue
		}

		// cancel context so that the previous instance of save_current_prices() stops execution
		cancel()
		// new context and cancellation func for new invocation of save_current_prices()
		ctx, cancel = context.WithCancel(context.TODO())

		go saveCurrentPrices(ctx, currentPrices)

		targetPrices, err := getTargetPrices()
		if err != nil {
			// TODO: log err
			continue
		}

		pricesUpdatePlan, err := compareCurrentVsTargetPrices(currentPrices, targetPrices)
		if err != nil {
			// TODO: log err
			continue
		}

		err = executePriceUpdatePlan(pricesUpdatePlan)
		if err != nil {
			// TODO: log err
			continue
		}
	}
}
