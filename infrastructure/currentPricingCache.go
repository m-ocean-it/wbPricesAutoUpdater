package infrastructure

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sync"
	"time"
)

func NewJsonCurrentPricingCache(filePath string) JsonCurrentPricingCache {
	var mu sync.Mutex
	return JsonCurrentPricingCache{filePath, &mu}
}

type JsonCurrentPricingCache struct {
	filePath string
	mu       *sync.Mutex
}

func (c *JsonCurrentPricingCache) GetPricing() ([]WbPricingItem, time.Duration, error) {
	c.mu.Lock()
	defer c.mu.Unlock()

	empty := []WbPricingItem{}

	fileContent, err := os.ReadFile(c.filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return empty, 0, fmt.Errorf("cache file does not exist: %w", err)
		}
		return empty, 0, err
	}

	var pricingItems []WbPricingItem
	err = json.Unmarshal(fileContent, &pricingItems)
	if err != nil {
		return empty, 0, err
	}

	fileStat, err := os.Stat(c.filePath)
	if err != nil {
		return empty, 0, err
	}

	cacheAge := time.Since(fileStat.ModTime())

	return pricingItems, cacheAge, nil
}

func (c *JsonCurrentPricingCache) SavePricing(data []WbPricingItem) error {
	jsonBytes, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshalling data to JSON: %w", err)
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	file, err := os.Create(c.filePath) // Create or open the file for writing (truncating if it already exists)
	if err != nil {
		return fmt.Errorf("error creating/opening file: %w", err)
	}
	defer file.Close()

	_, err = file.Write(jsonBytes)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}
