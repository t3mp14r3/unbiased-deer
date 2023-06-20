package domain

import (
    "errors"
    "net/http"

    "github.com/shopspring/decimal"
)

type User struct {
    ID          string          `db:"id" json:"id,omitempty"`
    Name        string          `db:"name" json:"name,omitempty"`
    Balance     decimal.Decimal `db:"balance" json:"balance,omitempty"`
}

type UserSimple struct {
    ID          string  `db:"id" json:"id,omitempty"`
    Name        string  `db:"name" json:"name,omitempty"`
    Balance     string  `db:"balance" json:"balance,omitempty"`
}

type Currency struct {
    ID          string          `db:"id" json:"id"`
    Symbol      string          `db:"symbol" json:"symbol"`
    Correlation decimal.Decimal `db:"correlation" json:"correlation"`
}

type CurrencySimple struct {
    ID          string  `db:"id" json:"id"`
    Symbol      string  `db:"symbol" json:"symbol"`
    Correlation string  `db:"correlation" json:"correlation"`
}

type RegisterRequest struct {
    Name    string  `json:"name"`
}

func (req *RegisterRequest) Validate() error {
    if len(req.Name) == 0 {
        return errors.New("name field is required")
    }

    return nil
}

type DepositRequest struct {
    Amount      string  `json:"amount"`
    Currency    string  `json:"currency"`
}

func (req *DepositRequest) Validate() error {
    if len(req.Amount) == 0 {
        return errors.New("amount is required")
    }

    amount, err := decimal.NewFromString(req.Amount)

    if err != nil {
        return err
    }

    if decimal.Zero.Cmp(amount) != -1 {
        return errors.New("amount must by positive")
    }
    
    if len(req.Currency) == 0 {
        return errors.New("currency is required")
    }

    return nil
}

type WithdrawRequest struct {
    Amount      string  `json:"amount"`
    Currency    string  `json:"currency"`
}

func (req *WithdrawRequest) Validate() error {
    if len(req.Amount) == 0 {
        return errors.New("amount is required")
    }

    amount, err := decimal.NewFromString(req.Amount)

    if err != nil {
        return err
    }

    if decimal.Zero.Cmp(amount) != -1 {
        return errors.New("amount must by positive")
    }
    
    if len(req.Currency) == 0 {
        return errors.New("currency is required")
    }

    return nil
}

type MessageType string

const (
    MessageRegister  MessageType = "register"
    MessageMe        MessageType = "me"
    MessageBalance   MessageType = "balance"
    MessageDeposit   MessageType = "deposit"
    MessageWithdraw  MessageType = "withdraw"
)

type ErrorType error

var (
    ErrorInternal   ErrorType = errors.New("something went wrong!")
    ErrorTimeout    ErrorType = errors.New("request timed out!")
    ErrorBadBody    ErrorType = errors.New("bad request body!")
    ErrorBadAmount  ErrorType = errors.New("invalid amount!")
)

type Message struct {
    Type    MessageType             `json:"type"`
    Payload map[string]interface{}  `json:"payload"`
}

type Response struct {
    Err     ErrorType
    Data    map[string]interface{}
}

func (r Response) Status() int {
    switch r.Err {
        case ErrorInternal:
            return http.StatusInternalServerError
        case ErrorTimeout:
            return http.StatusRequestTimeout
        case ErrorBadBody:
            return http.StatusBadRequest
        case ErrorBadAmount:
            return http.StatusBadRequest
        case nil:
            return http.StatusOK
        default:
            return http.StatusInternalServerError
    }
}
