package pricingTableHandler

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type PricingTableHandler struct {
	Db *sqlx.DB
}

func (p PricingTableHandler) WriteOverrideMpPrice(wbSku string, price uint) error {
	_, err := p.Db.NamedExec(`
			UPDATE pricing
			SET
				override_mp_price = :override_mp_price,
				revised_at = now()
			WHERE wb_sku = :wb_sku
		`,
		Row{
			WbSku:           wbSku,
			OverrideMpPrice: price,
		})
	if err != nil {
		return fmt.Errorf("could not write override-price to database: %w", err)
	}

	return nil
}

func (p PricingTableHandler) RemoveOverrideMpPrice(wbSku string) error {
	_, err := p.Db.NamedExec(`
			UPDATE pricing
			SET
				override_mp_price = NULL,
				revised_at = now()
			WHERE wb_sku = :wb_sku
		`,
		Row{
			WbSku: wbSku,
		})

	if err != nil {
		return fmt.Errorf("could not nullify override-price in the database: %w", err)
	}

	return nil
}

func (p PricingTableHandler) JustRevise(wbSku string) error {
	_, err := p.Db.NamedExec(`
			UPDATE pricing
			SET
				revised_at = now()
			WHERE wb_sku = :wb_sku
		`,
		Row{
			WbSku: wbSku,
		})

	if err != nil {
		return fmt.Errorf("could not revise wb sku pricing in the database: %w", err)
	}

	return nil
}
