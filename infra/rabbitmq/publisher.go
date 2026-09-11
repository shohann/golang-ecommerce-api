package rabbitmq

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shohann/golang-ecommerce-api/order"
)

type Publisher struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
}

func NewPublisher(conn *amqp.Connection, queueName string) (*Publisher, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	if _, err := declareQueue(ch, queueName); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("rabbitmq declare queue: %w", err)
	}

	return &Publisher{
		conn:      conn,
		ch:        ch,
		queueName: queueName,
	}, nil
}

func (p *Publisher) PublishOrderPlaced(event order.OrderPlacedEvent) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("marshal order placed event: %w", err)
	}

	err = p.ch.Publish(
		"",
		p.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			DeliveryMode: amqp.Persistent,
			Body:         body,
		},
	)
	if err != nil {
		return fmt.Errorf("publish order placed: %w", err)
	}

	return nil
}

func (p *Publisher) Close() error {
	if p.ch != nil {
		if err := p.ch.Close(); err != nil {
			return err
		}
	}
	return nil
}
