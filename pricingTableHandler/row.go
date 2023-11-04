package pricingTableHandler

type Row struct {
	WbSku           string `db:"wb_sku"`
	OverrideMpPrice uint   `db:"override_mp_price"`
}
