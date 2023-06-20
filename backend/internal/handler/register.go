package handler

import (
	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
)

func (h *Handler) Register(message domain.Message) domain.Response {
    token, err := h.repo.CreateUser(message.Payload["name"].(string))

    var response domain.Response

    if err != nil {
        response.Err = domain.ErrorInternal
    } else {
        response.Data = map[string]interface{}{"token": token}
    }

    return response
}
