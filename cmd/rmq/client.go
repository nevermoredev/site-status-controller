package rmq

import (
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"log"
	RmqProto "zeithub.com/site-status-controller/pkg/config/protobuf"
)

func Send(body *RmqProto.BotJobResponse) {
	conn, err := amqp.Dial("amqp://test:password@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	response, err := proto.Marshal(body)

	q, err := ch.QueueDeclare(
		"spider-log", // name
		false,         // durable
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
			Body:        response,
		})
	log.Print(response)
	failOnError(err, "Failed to publish a message")

}
