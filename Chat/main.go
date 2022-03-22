package main

import (
	"encoding/json"
	"goChallenge/chat/broker"
	"goChallenge/chat/config"
	"goChallenge/chat/db"
	"goChallenge/chat/server"
	"log"
)

func main() {

	server.RoomsMessages = make(map[string][]string)

	config.Load()

	// MySQL Connection
	db.Init()
	defer db.Close()

	// RabbitMQ connection
	amqp, err := broker.Connect()
	if err != nil {
		log.Fatal(err)
	}
	defer amqp.Close()
	ch, err := broker.OpenChannel()
	if err != nil {
		return
	}
	defer ch.Close()

	server := server.NewServer()
	go server.Run()

	msgs := broker.ReceiveMessageDeliveryChannel()

	go func() {
		for d := range msgs {
			var response broker.ClientMessage
			json.Unmarshal(d.Body, &response)

			hubs := *server.GetHubs()
			hub := hubs[response.HubName]
			hub.SendTo(response.Message, response.ClientRemoteAddress)
		}
	}()

	server.Start()
}
