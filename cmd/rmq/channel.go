package rmq

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/zeithub/site-status-controller/cmd/checker"
	RmqProto "github.com/zeithub/site-status-controller/pkg/config/protobuf"
	"google.golang.org/protobuf/proto"
)

// func failOnError(err error, msg string) {
// 	if err != nil {
// 		log.Fatalf("%s: %s", msg, err)
// 	}
// }

func Listen(rmqConn *amqp.Connection) {

	ch1, errChan1 := rmqConn.Channel()
	if errChan1 != nil {
		log.Fatalf("Error (channel.Listen.Chan1): %v", errChan1)
	}
	defer ch1.Close()

	ch2, errChan2 := rmqConn.Channel()
	if errChan2 != nil {
		log.Fatalf("Error (channel.Listen.Chan2): %v", errChan2)
	}
	defer ch2.Close()

	q1, errQueue1 := ch1.QueueDeclare(
		"spider-feed", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	if errQueue1 != nil {
		log.Fatalf("Error (channel.Listen.Queue1): %v", errQueue1)
	}

	q2, errQueue2 := ch2.QueueDeclare(
		"spider-log", // name
		false,        // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait (wait time for processing)
		nil,          // arguments
	)
	if errQueue2 != nil {
		log.Fatalf("Error (channel.Listen.Queue2): %v", errQueue2)
	}

	args := make(amqp.Table)
	args["x-queue-mode"] = "lazy"

	msgs, errConsume := ch1.Consume(
		q1.Name,       // queue
		"site-status", // consumer
		true,          // auto-ack
		false,         // exclusive
		false,         // no-local
		false,         // no-wait
		args,          // args
	)
	if errConsume != nil {
		log.Fatalf("Error (channel.Listen.Consumer): %v", errConsume)
	}

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			p := &RmqProto.BotJob{}
			if err := proto.Unmarshal(d.Body, p); err != nil {
				log.Fatalln("Failed to parse Person:", err)
			}
			// log.Printf("%s",p.PageUrl)
			pageInfo := checker.TestSite(p.PageId, p.PageUrl, p.Title, p.Action)
			Send(ch2, q2, pageInfo)
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}
