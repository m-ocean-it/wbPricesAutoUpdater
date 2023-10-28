-- +goose Up
-- +goose StatementBegin
CREATE TABLE public.pricing (
	wb_sku varchar NOT NULL,
	mp_price smallint NULL,
	override_mp_price smallint NULL,
	current_wb_price smallint NULL,
	current_wb_discount smallint NULL,
	note text null,
	CONSTRAINT pricing_pk PRIMARY KEY (wb_sku)
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE public.pricing;
-- +goose StatementEnd
