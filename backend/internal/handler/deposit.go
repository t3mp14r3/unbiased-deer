package handler

import (
	"go.uber.org/zap"

	"github.com/shopspring/decimal"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
)

func (h *Handler) Deposit(message domain.Message) domain.Response {
    user, err := h.repo.GetUser(message.Payload["id"].(string))

    if err != nil {
        return domain.Response{Err: domain.ErrorInternal}
    }
    
    currency, err := h.repo.GetCurrency(message.Payload["currency"].(string))
    
    if err != nil {
        return domain.Response{Err: domain.ErrorInternal}
    }

    deposit, err := decimal.NewFromString(message.Payload["amount"].(string))
    
    if err != nil {
        h.logger.Error("failed to parse deposit amount", zap.Error(err))
        return domain.Response{Err: domain.ErrorInternal}
    }

    deposit = deposit.Div(currency.Correlation)
    user.Balance = user.Balance.Add(deposit)
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
