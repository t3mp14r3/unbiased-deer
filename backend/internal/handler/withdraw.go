package handler

import (
	"go.uber.org/zap"

	"github.com/shopspring/decimal"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
)

func (h *Handler) Withdraw(message domain.Message) domain.Response {
    user, err := h.repo.GetUser(message.Payload["id"].(string))

    if err != nil {
        return domain.Response{Err: domain.ErrorInternal}
    }
    
    currency, err := h.repo.GetCurrency(message.Payload["currency"].(string))
    
    if err != nil {
        return domain.Response{Err: domain.ErrorInternal}
    }

    withdraw, err := decimal.NewFromString(message.Payload["amount"].(string))
    
    if err != nil {
        h.logger.Error("failed to parse withdraw amount", zap.Error(err))
        return domain.Response{Err: domain.ErrorInternal}
    }

    withdraw = withdraw.Div(currency.Correlation)

    if withdraw.Cmp(user.Balance) == 1 {
        h.logger.Error("withdraw amount is too large", zap.Error(err))
        return domain.Response{Err: domain.ErrorBadAmount}
    }

    user.Balance = user.Balance.Sub(withdraw)
    displayBalance := user.Balance.Mul(currency.Correlation).Round(2).String()

    err = h.repo.UpdateUser(user)
    
    var response domain.Response

    if err != nil {
        return domain.Response{Err: domain.ErrorInternal}
    } else {
        response.Data = map[string]interface{}{
            "name": user.Name,
            "balance": displayBalance,
            "currency": currency.ID,
        }
    }

    return response
}
