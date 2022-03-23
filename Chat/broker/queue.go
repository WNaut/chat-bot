package broker

import (
	"encoding/json"
	"fmt"
	"goChallenge/chat/config"
	"log"

	"github.com/streadway/amqp"
)

type ClientMessage struct {
	HubName             string `json:"hubName"`
	ClientRemoteAddress string `json:"clientRemoteAddress"`
	Message             string `json:"message"`
}

var (
	Conn                            *amqp.Connection
	Channel                         *amqp.Channel
	ClientQueueName, StooqQueueName string
)

// start the connection with rabbit
func Connect() (*amqp.Connection, error) {
	url := fmt.Sprintf("amqp://%s:%s@%s:%s/",
		config.RabbitConfig.User,
		config.RabbitConfig.Pwd,
		config.RabbitConfig.Host,
		config.RabbitConfig.Port)
	conn, err := amqp.Dial(url)

	if err != nil {
		failOnError(err, "Failed to connect to RabbitMQ")
		return nil, err
	}

	Conn = conn

	ClientQueueName = config.RabbitConfig.ClientQueue
	StooqQueueName = config.RabbitConfig.StooqQueue

	return conn, nil
}

func OpenChannel() (*amqp.Channel, error) {
	ch, err := Conn.Channel()
	if err != nil {
		failOnError(err, "Failed to open channel")
		return nil, err
	}

	Channel = ch

	return ch, nil
}

func SendMessage(message *ClientMessage) {
	ch, err := Conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		ClientQueueName, // name
		false,           // durable
		false,           // delete when unused
		false,           // exclusive
		false,           // no-wait
		nil,             // arguments
	)
	failOnError(err, "Failed to declare a queue")

	json, err := json.Marshal(message)
	if err != nil {
		failOnError(err, "Failed to parse body message")
	}

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        json,
		})
	failOnError(err, "Failed to publish a message")

	fmt.Printf("Message sent: %s\n", json)
}

func ReceiveMessageDeliveryChannel() <-chan amqp.Delivery {
	q, err := Channel.QueueDeclare(
		StooqQueueName, // name
		false,          // durable
		false,          // delete when unused
		false,          // exclusive
		false,          // no-wait
		nil,            // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := Channel.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	return msgs
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
