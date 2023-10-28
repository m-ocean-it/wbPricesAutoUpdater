package infrastructure

import (
	"encoding/json"
	"fmt"
	"time"
)

type JsonCurrentPricingCache struct {
	FileHandler ConcurrentFileHandler
}

func (c *JsonCurrentPricingCache) GetPricing() ([]WbPricingItem, time.Duration, error) {
	empty := []WbPricingItem{}

	fileContent, fileStat, err := c.FileHandler.Read()
	if err != nil {
		return empty, 0, fmt.Errorf("could not read file: %w", err)
	}

	var pricingItems []WbPricingItem
	err = json.Unmarshal(fileContent, &pricingItems)
	if err != nil {
		return empty, 0, fmt.Errorf("could not unmarshal data to pricing items: %w", err)
	}

	cacheAge := time.Since(fileStat.ModTime())

	return pricingItems, cacheAge, nil
}

func (c *JsonCurrentPricingCache) SavePricing(data []WbPricingItem) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling pricing data to JSON: %w", err)
	}

	err = c.FileHandler.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("error writing pricing to file: %w", err)
	}

	return nil
}
