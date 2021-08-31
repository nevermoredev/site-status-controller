package rmq

import (
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
	"log"
	"zeithub.com/site-status-controller/api/config"
	"zeithub.com/site-status-controller/build"
	"zeithub.com/site-status-controller/cmd/checker"
)

type Message struct {
	PageId string `json:PageId`
	Url    string `json:Url`
	Action string `json:Action`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

func Listen() {
	configRMQ := config.GetConfig()
	args := make(amqp.Table)
	args["x-queue-mode"] = "lazy"
	conn, err := amqp.Dial("amqp://test:password@" + configRMQ.Host + ":" + configRMQ.Port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"botjobs", // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name,    // queue
		"service", // consumer
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
			p := &pb.BotJob{}
			if err := proto.Unmarshal(d.Body, p); err != nil {
				log.Fatalln("Failed to parse Person:", err)
			}
			//log.Printf("%s", err)
			var titleNow = checker.CheckTitle(p.PageUrl)
			//log.Printf("%s %s %d %s", SiteInfo.Url, SiteInfo.TitleNow, SiteInfo.Status, SiteInfo.Hash)
			go Send(titleNow)

		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
