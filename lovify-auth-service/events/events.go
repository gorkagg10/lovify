package events

import (
	"context"
	"fmt"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
	"log/slog"
)

const (
	CreateProfileConsumer = "createProfileConsumer"

	Workqueue = "workqueue"
)

type NatsHandler struct {
	jetStream jetstream.JetStream
	natsConn  *nats.Conn
}

func NewNatsHandler(natsURL string) (*NatsHandler, error) {
	natsConn, err := nats.Connect(natsURL)
	if err != nil {
		return nil, err
	}
	jetStream, err := jetstream.New(natsConn)
	if err != nil {
		return nil, err
	}
	return &NatsHandler{
		natsConn:  natsConn,
		jetStream: jetStream,
	}, nil
}

type Consumer struct {
	name   string
	stream string
	msg    func(msg jetstream.Msg)
}

func NewConsumer(name string, stream string, msg func(msg jetstream.Msg)) Consumer {
	return Consumer{
		name:   name,
		stream: stream,
		msg:    msg,
	}
}

func (h NatsHandler) Consume(ctx context.Context, stream, consumer string, msg func(msg jetstream.Msg)) error {
	con, err := h.jetStream.Consumer(ctx, stream, consumer)
	if err != nil {
		return err
	}
	_, err = con.Consume(msg)
	return err
}

func Listen(ctx context.Context, natsHandler *NatsHandler, consumers []Consumer) error {
	errors := []error{}
	for _, consumer := range consumers {
		if err := natsHandler.Consume(ctx, consumer.stream, consumer.name, consumer.msg); err != nil {
			slog.Info("err", "error", err)
			errors = append(errors, err)
		}
		slog.Info("consumer created", slog.String("consumer", consumer.name))
	}
	if len(errors) > 0 {
		return fmt.Errorf("listen failed: %v", errors)
	}
	<-ctx.Done()
	if ctx.Err() == context.DeadlineExceeded {
		return nil
	}
	return ctx.Err()
}
