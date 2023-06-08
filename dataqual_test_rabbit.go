package dataqual_wasm

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/streamdal/dataqual"
	amqp "github.com/streamdal/rabbitmq-amqp091-go"
)

type Rabbit struct {
	Channel    *amqp.Channel
	Connection *amqp.Connection
}

func setupRabbit() (*Rabbit, error) {
	ac, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		return nil, errors.Wrap(err, "unable to dial rabbit")
	}

	ch, err := ac.Channel()
	if err != nil {
		return nil, errors.Wrap(err, "unable to instantiate channel")
	}

	if err := ch.Qos(1, 0, false); err != nil {
		return nil, errors.Wrap(err, "unable to set qos policy")
	}

	if err := ch.ExchangeDeclare(AmqpExchange, amqp.ExchangeTopic, false, true, false, false, nil); err != nil {
		return nil, errors.Wrap(err, "unable to declare exchange")
	}

	if _, err := ch.QueueDeclare(AmqpQueue, true, false, false, false, nil); err != nil {
		return nil, errors.Wrap(err, "unable to declare queue")
	}

	if err := ch.QueueBind(AmqpQueue, AmqpRoutingKey, AmqpExchange, false, nil); err != nil {
		return nil, errors.Wrap(err, "unable to bind queue")
	}

	return &Rabbit{
		Channel:    ch,
		Connection: ac,
	}, nil
}

func (r *Rabbit) Cleanup() {
	defer r.Connection.Close()
	defer r.Channel.Close()

	if _, err := r.Channel.QueueDelete(AmqpQueue, false, false, false); err != nil {
		panic(err)
	}

	if err := r.Channel.ExchangeDelete(AmqpExchange, false, false); err != nil {
		panic(err)
	}
}

func (r *Rabbit) Produce(data []byte) {
	err := r.Channel.PublishWithContext(context.Background(), AmqpExchange, "dqtest", false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		Body:         data,
		AppId:        "testing",
	})
	if err != nil {
		if !strings.Contains(err.Error(), dataqual.ErrMessageDropped.Error()) {
			panic("unable to publish message: " + err.Error())
		}
	}
}
