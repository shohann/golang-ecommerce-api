package rabbitmq

import (
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/shohann/golang-ecommerce-api/order"
)

type OrderPlacedHandler func(event order.OrderPlacedEvent) error

type Consumer struct {
	conn      *amqp.Connection
	ch        *amqp.Channel
	queueName string
}

func NewConsumer(conn *amqp.Connection, queueName string) (*Consumer, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq channel: %w", err)
	}

	if _, err := declareQueue(ch, queueName); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("rabbitmq declare queue: %w", err)
	}

	if err := ch.Qos(1, 0, false); err != nil {
		_ = ch.Close()
		return nil, fmt.Errorf("rabbitmq qos: %w", err)
	}

	return &Consumer{
		conn:      conn,
		ch:        ch,
		queueName: queueName,
	}, nil
}

func (c *Consumer) ConsumeOrderPlaced(handler OrderPlacedHandler) error {
	deliveries, err := c.ch.Consume(
		c.queueName,
		"",
		false, // manual ack
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return fmt.Errorf("rabbitmq consume: %w", err)
	}

	for d := range deliveries {
		var event order.OrderPlacedEvent
		if err := json.Unmarshal(d.Body, &event); err != nil {
			_ = d.Nack(false, false)
			fmt.Println("failed to unmarshal order placed event:", err)
			continue
		}

		if err := handler(event); err != nil {
			_ = d.Nack(false, true)
			fmt.Println("order placed handler error:", err)
			continue
		}

		if err := d.Ack(false); err != nil {
			fmt.Println("failed to ack message:", err)
		}
	}

	return nil
}

func (c *Consumer) Close() error {
	if c.ch != nil {
		if err := c.ch.Close(); err != nil {
			return err
		}
	}
	return nil
}
