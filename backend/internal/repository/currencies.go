package repository

import (
	"github.com/shopspring/decimal"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
	"go.uber.org/zap"
)

func (r *RepoClient) GetCurrency(id string) (domain.Currency, error) {
    query := `SELECT id, symbol, correlation FROM currencies WHERE LOWER(id) = LOWER($1);`

    var tmp domain.CurrencySimple
    err := r.pc.GetContext(r.ctx, &tmp, query, id)
    
    if err != nil {
        r.logger.Error("failed to query currency record", zap.Error(err))
        return domain.Currency{}, err
    }

    var currency domain.Currency
    currency.ID = tmp.ID
    currency.Symbol = tmp.Symbol
    
    currency.Correlation, err = decimal.NewFromString(tmp.Correlation)

    if err != nil {
        r.logger.Error("failed to parse correlation value", zap.Error(err))
        return domain.Currency{}, err
    }

    return currency, nil
}
