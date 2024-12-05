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
	fmt.Println("Starting Peril client...")
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
	fmt.Println("Peril game client connected to RabbitMQ!")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}
	_, queue, err := pubsub.DeclareAndBind(
		conn,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+username,
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)
	gamestate := gamelogic.NewGameState(username)
	for {
		words := gamelogic.GetInput()
		if words[0] == "spawn" {
			err := gamestate.CommandSpawn(words)
			if err != nil {
				log.Printf("could not spawn: %v", err)
			}

		} else if words[0] == "move" {
			_, err := gamestate.CommandMove(words)
			if err != nil {
				log.Printf("could not move: %v", err)
			}
		} else if words[0] == "status" {
			gamestate.CommandStatus()
		} else if words[0] == "help" {
			gamelogic.PrintClientHelp()
		} else if words[0] == "spam" {
			fmt.Println("Spamming not allowed yet!")
		} else if words[0] == "quit" {
			gamelogic.PrintQuit()
			break
		} else {
			fmt.Println("Unknown command...")
		}
	}
}
