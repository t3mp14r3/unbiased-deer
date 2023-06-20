package handler

import (
	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
)

func (h *Handler) Me(message domain.Message) domain.Response {
    user, err := h.repo.GetUser(message.Payload["id"].(string))
    
    if err != nil {
        return domain.Response{Err: domain.ErrorInternal}
    }
    
    currency, err := h.repo.GetCurrency(message.Payload["currency"].(string))
    
    if err != nil {
        return domain.Response{Err: domain.ErrorInternal}
    }

    user.Balance = user.Balance.Mul(currency.Correlation)

    var response domain.Response

    response.Data = map[string]interface{}{
        "name": user.Name,
        "balance": user.Balance.Round(2).String(),
        "currency": currency.ID,
    }

    return response
}
