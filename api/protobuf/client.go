package rmq

import (
	"github.com/streadway/amqp"
	"log"
)

func Send(body string) {
	conn, err := amqp.Dial("amqp://test:password@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"botresults", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait (wait time for processing)
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	err = ch.Publish(
		"",     // exchange
		q.Name, // routing key
		false,  // mandatory
		false,  // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(body),
		})
	log.Printf(" [x] Sent %s", body)
	failOnError(err, "Failed to publish a message")

}
