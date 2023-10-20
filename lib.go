package main

import "context"

type price uint16

type task struct {
	productId uint64
	setPrice  price
}

func getCurrentPrices() ([]price, error) {
	return []price{}, nil
}

func saveCurrentPrices(ctx context.Context, prices []price) error {
	return nil
}

func getTargetPrices() ([]price, error) {
	return []price{}, nil
}

func compareCurrentVsTargetPrices(current []price, target []price) ([]task, error) {
	return []task{}, nil
}

func executeTasks(tasks []task) error {
	return nil
}
