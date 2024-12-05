package main

import (
	"fmt"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
	"log"
)

func main() {
	fmt.Println("Starting Peril server...")
	const connString = "amqp://guest:guest@localhost:5672/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer func(conn *amqp.Connection) {
		err := conn.Close()
		if err != nil {
			log.Fatalf("Failed to close connection to RabbitMQ server: %v", err)
		}
	}(conn)
	fmt.Println("Peril game server connected to RabbitMQ!")
	rmqChan, err := conn.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}
	_, _, err = pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routing.GameLogSlug+".*",
		pubsub.SimpleQueueDurable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	gamelogic.PrintServerHelp()
	for {
		words := gamelogic.GetInput()
		if words[0] == "pause" {
			fmt.Println("Pausing game...")
			err = pubsub.PublishJSON(
				rmqChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: true},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		} else if words[0] == "resume" {
			fmt.Println("Resuming game...")
			err = pubsub.PublishJSON(
				rmqChan,
				routing.ExchangePerilDirect,
				routing.PauseKey,
				routing.PlayingState{IsPaused: false},
			)
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}
		} else if words[0] == "quit" {
			fmt.Println("Exiting Peril server...")
			break
		} else {
			fmt.Println("Unknown command...")
		}
	}
}
