-- +goose Up
-- +goose StatementBegin
ALTER TABLE pricing
ADD COLUMN revised_at timestamp with time zone NOT NULL DEFAULT now();

ALTER TABLE pricing ALTER COLUMN revised_at DROP DEFAULT;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE pricing
DROP COLUMN revised_at;
-- +goose StatementEnd
