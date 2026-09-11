package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/shohann/golang-ecommerce-api/config"
	"github.com/shohann/golang-ecommerce-api/infra/rabbitmq"
	"github.com/shohann/golang-ecommerce-api/order"
)

func main() {
	cnf := config.GetConfig()

	conn, err := rabbitmq.Connect(cnf.RabbitMQURL)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer conn.Close()

	consumer, err := rabbitmq.NewConsumer(conn, cnf.RabbitMQOrderQueue)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer consumer.Close()

	fmt.Println("order worker listening on queue:", cnf.RabbitMQOrderQueue)

	go func() {
		err := consumer.ConsumeOrderPlaced(func(event order.OrderPlacedEvent) error {
			fmt.Printf("order placed: %#v\n", event)
			return nil
		})
		if err != nil {
			fmt.Println("consumer stopped:", err)
			os.Exit(1)
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
	<-sigCh
	fmt.Println("shutting down worker")
}
