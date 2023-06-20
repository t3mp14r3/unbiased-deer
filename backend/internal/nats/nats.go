package nats

import (
	"context"
	"encoding/json"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"go.uber.org/zap"

	"github.com/t3mp14r3/unbiased-deer/backend/internal/config"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/domain"
	"github.com/t3mp14r3/unbiased-deer/backend/internal/handler"
)

type NatsClient struct {
    conn    *nats.Conn
    logger  *zap.Logger
    handler *handler.Handler
}

func New(natsConfig *config.NatsConfig, logger *zap.Logger, handler *handler.Handler) *NatsClient {
    conn, err := nats.Connect(natsConfig.Url)

    if err != nil {
        log.Fatalln("failed to initialize nats connection! err:", err)
    }

    return &NatsClient{
        conn:       conn,
        logger:     logger,
        handler:    handler,
    }
}

func (n *NatsClient) Send(message domain.Message) domain.Response {
    bytes, err := json.Marshal(message)

    if err != nil {
        n.logger.Error("failed to marshal the request", zap.Error(err))
        return domain.Response{Err: domain.ErrorInternal}
    }

    resp, err := n.conn.Request("event", bytes, time.Duration(time.Second * 5))

    if err != nil {
        n.logger.Error("failed to send the request", zap.Error(err))
        return domain.Response{Err: domain.ErrorTimeout}
    }

    var result domain.Response

    err = json.Unmarshal(resp.Data, &result)
    
    if err != nil {
        n.logger.Error("failed to unmarshal the response", zap.Error(err))
        return domain.Response{Err: domain.ErrorInternal}
    }
        
    return result
}

func (n *NatsClient) Subscribe(ctx context.Context) error {
    _, err := n.conn.Subscribe("event", n.reply)

    if err != nil {
        n.logger.Error("failed to subscribe to an event", zap.Error(err))
        return err
    }

    <-ctx.Done()

    n.conn.Drain()

    loop:for {
        if n.conn.IsClosed() {
            break loop
        }
    }

    return nil
}

func (n *NatsClient) reply(msg *nats.Msg) {
    var response domain.Response
    var message domain.Message

    err := json.Unmarshal(msg.Data, &message)
    
    if err != nil {
        n.logger.Error("failed to unmarshal the response", zap.Error(err))
        response = domain.Response{Err: domain.ErrorInternal}
        return
    }

    switch message.Type {
        case domain.MessageRegister:
            response = n.handler.Register(message)
        case domain.MessageMe:
            response = n.handler.Me(message)
        case domain.MessageDeposit:
            response = n.handler.Deposit(message)
        case domain.MessageWithdraw:
            response = n.handler.Withdraw(message)
    }
    
    bytes, err := json.Marshal(response)

    if err != nil {
        n.logger.Error("failed to marshal the response", zap.Error(err))
        return
    }

    msg.Respond(bytes)
}
