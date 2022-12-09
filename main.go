package main

import (
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	conn := con()
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Faild to open channel")
	defer ch.Close()

	/* create queue*/
	q, err := ch.QueueDeclare(
		"queue_green",
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	failOnError(err, "Faild to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			/* manual ack */
			//d.Ack(false)
		}
	}()

	log.Printf(" [*] Waiting for messages, To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func con() *amqp.Connection {
	conn, err := amqp.Dial("amqp://rabbit:password@localhost:5672")
	failOnError(err, "Faild to connect to RabbitMQ")

	return conn
}
