package main

import (
	"encoding/json"
	"fmt"
	"goChallenge/chatbot/queue"
	"goChallenge/chatbot/stockService"
	"log"
)

func main() {
	amqp, err := queue.Connect()
	if err != nil {
		return
	}
	defer amqp.Close()
	ch, err := queue.OpenChannel()
	if err != nil {
		return
	}
	defer ch.Close()

	msgs := queue.ReceiveMessageDeliveryChannel()
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var cm queue.ClientMessage
			json.Unmarshal(d.Body, &cm)
			fmt.Println(d, cm)
			message, err := stockService.GetStockQuote(cm.Message)
			response := queue.ClientMessage{HubName: cm.HubName, ClientRemoteAddress: cm.ClientRemoteAddress, Message: message}
			if err != nil {
				//log.Fatal(err)
				return
			}

			clientMessage := &queue.ClientMessage{HubName: response.HubName, ClientRemoteAddress: response.ClientRemoteAddress, Message: response.Message}
			queue.SendMessage(clientMessage)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
