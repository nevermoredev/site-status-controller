package rmq

import (
	"log"

	"github.com/streadway/amqp"
	RmqProto "github.com/zeithub/site-status-controller/pkg/config/protobuf"
	"google.golang.org/protobuf/proto"
)

func Send(ch2 *amqp.Channel, queue2 amqp.Queue, body *RmqProto.BotJobResponse) {

	response, err := proto.Marshal(body)

	err = ch2.Publish(
		"",          // exchange
		queue2.Name, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        response,
		},
	)
	if err != nil {
		log.Printf("Error (client.Send): %v", err)
		return
	}

}
