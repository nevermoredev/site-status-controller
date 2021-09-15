package main

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/zeithub/site-status-controller/cmd/rmq"
)

func main() {

	rmqConn, errRmq := amqp.Dial("amqp://metrix:digitAlks~256@localhost:5672/")
	if errRmq != nil {
		log.Fatalf("Error connecting to RabbitMQ: %v", errRmq)
	}
	defer rmqConn.Close()

	rmq.Listen(rmqConn)

}
