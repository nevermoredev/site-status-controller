package rmq

import (
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"log"
	"zeithub.com/site-status-controller/cmd/checker"
	"zeithub.com/site-status-controller/pkg/config/protobuf"
	"zeithub.com/site-status-controller/pkg/config/rmq"
)



func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Listen() {
	configRMQ := rmq.GetConfig()
	args := make(amqp.Table)
	args["x-queue-mode"] = "lazy"
	conn, err := amqp.Dial("amqp://test:password@" + configRMQ.Host + ":" + configRMQ.Port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"spider-feed", // name
		false,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,    // queue
		"site-status", // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no-local
		false,     // no-wait
		args,      // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			p := &RmqProto.BotJobResponse{}
			if err := proto.Unmarshal(d.Body, p); err != nil {
				log.Fatalln("Failed to parse Person:", err)
			}
			log.Printf("%s",p.PageUrl)
			 	pageInfo := checker.TestSite(p.PageId,p.PageUrl)
			Send(pageInfo)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
