package repository

import (
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"github.com/shopspring/decimal"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
	"go.uber.org/zap"
)

func (r *RepoClient) CreateUser(name string) (string, error) {
    query := `INSERT INTO users(name) VALUES($1) RETURNING id;`

    var id string

    err := r.pc.GetContext(r.ctx, &id, query, name)

    if err != nil {
        r.logger.Error("failed to create new user record", zap.Error(err))
    }

    var token string

    loop:for {
        token = uuid.NewString()

        if err := r.rc.Get(r.ctx, token).Err(); err == redis.Nil {
            break loop
        }
    }

    err = r.rc.Set(r.ctx, token, id, 0).Err()

    if err != nil {
        r.logger.Error("failed to create new user session record", zap.Error(err))
    }

    return token, err
}

func (r *RepoClient) GetUser(id string) (domain.User, error) {
    query := `SELECT 
        id,
        name,
        balance 
        FROM users 
        WHERE id = $1;`

    var tmp domain.UserSimple
    err := r.pc.GetContext(r.ctx, &tmp, query, id)

    if err != nil {
        r.logger.Error("failed to query user record", zap.Error(err))
        return domain.User{}, err
    }

    var user domain.User
    user.ID = tmp.ID
    user.Name = tmp.Name
    
    user.Balance, err = decimal.NewFromString(tmp.Balance)

    if err != nil {
        r.logger.Error("failed to parse correlation value", zap.Error(err))
        return domain.User{}, err
    }

    return user, err
}

func (r *RepoClient) UpdateUser(user domain.User) error {
    query := `UPDATE users SET 
        name = $1, 
        balance = $2 
        WHERE id = $3;`

    _, err := r.pc.ExecContext(r.ctx, query, user.Name, user.Balance.String(), user.ID)

    if err != nil {
        r.logger.Error("failed to update user record", zap.Error(err))
    }

    return err
}
