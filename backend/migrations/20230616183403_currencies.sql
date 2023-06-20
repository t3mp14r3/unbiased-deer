-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS currencies (
    id                  TEXT PRIMARY KEY,
    symbol              TEXT NOT NULL,
    correlation         TEXT NOT NULL
);
INSERT INTO currencies(id, symbol, correlation) VALUES('USD', '$', '0.012') ON CONFLICT DO NOTHING;
INSERT INTO currencies(id, symbol, correlation) VALUES('EUR', '€', '0.011') ON CONFLICT DO NOTHING;
INSERT INTO currencies(id, symbol, correlation) VALUES('RUB', '₽', '1') ON CONFLICT DO NOTHING;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS currencies;
-- +goose StatementEnd
