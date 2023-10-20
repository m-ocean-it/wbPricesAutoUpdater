package domain

type ProductId uint64
type Price uint16
type Discount uint8

type PricePair struct {
	Price    Price
	Discount Discount
}
type CatalogPricing map[ProductId]PricePair
