package main

import (
	"fmt"
	"log"
	"github.com/streadway/amqp"
	"github.com/go-logging-micro-service/logger"
	"encoding/json"
	"github.com/joho/godotenv"
)

func failOnError(err error, msg string) {
	if err != nil {
		// Exit the program.
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}


	// 'rabbitmq-server' is the network reference we have to the broker,
	conn, err := amqp.Dial("amqp://guest:guest@rabbitmq_api:5672/")
	failOnError(err, "Error connecting to the broker")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	exchangeName := "api_services"
	bindingKey   := "service.logging.*"

	// Create the exchange if it doesn't already exist.
	err = ch.ExchangeDeclare(
			exchangeName, 	// name
			"topic",  		// type
			true,         	// durable
			false,
			false,
			false,
			nil,
	)
	failOnError(err, "Error creating the exchange")

	q, err := ch.QueueDeclare(
			"service_queue",    // name - empty means a random, unique name will be assigned
			true,  // durable
			false, // delete when unused
			false,
			false,
			nil,
	)
	failOnError(err, "Error creating the queue")

	// Bind the queue to the exchange based on a string pattern (binding key).
	err = ch.QueueBind(
			q.Name,       // queue name
			bindingKey,   // binding key
			exchangeName, // exchange
			false,
			nil,
	)
	failOnError(err, "Error binding the queue")

	// Subscribe to the queue.
	msgs, err := ch.Consume(
			q.Name, // queue
			"",     // consumer id - empty means a random, unique id will be assigned
			false,  // auto acknowledgement of message delivery
			false,
			false,
			false,
			nil,
	)
	failOnError(err, "Failed to register as a consumer")
	forever := make(chan bool)

	go func() {
		for d := range msgs {
			log.Printf("Received message: %s", d.Body)

			m := map[string]interface{}{}
			err := json.Unmarshal([]byte(string(d.Body)), &m)
			if err != nil {
			    panic(err)
			}

			if m["log_data"] != nil {
				logger.PutLog(m["log_data"])
			}

			d.Ack(false)
		}
	}()

	fmt.Println("Service listening for events...")

	// Block until 'forever' receives a value
	// No value will be sent in order continuosly listen for new messages
	<-forever
}
