package rmq

import (
	uuid "github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
	"log"
	"time"
)

func (rcl *RabbitClient) channel(isRec, recreate bool) (*amqp.Channel, error) {
	if recreate {
		if isRec {
			rcl.recChan = nil
		} else {
			rcl.sendChan = nil
		}
	}
	if isRec && rcl.recConn == nil {
		rcl.recChan = nil
	}
	if !isRec && rcl.sendConn == nil {
		rcl.recChan = nil
	}
	if isRec && rcl.recChan != nil {
		return rcl.recChan, nil
	} else if !isRec && rcl.sendChan != nil {
		return rcl.sendChan, nil
	}
	for {
		_, err := rcl.connect(isRec, recreate)
		if err == nil {
			break
		}
	}
	var err error
	if isRec {
		rcl.recChan, err = rcl.recConn.Channel()
	} else {
		rcl.sendChan, err = rcl.sendConn.Channel()
	}
	if err != nil {
		log.Println("--- could not create channel ---")
		time.Sleep(1 * time.Second)
		return nil, err
	}
	if isRec {
		return rcl.recChan, err
	} else {
		return rcl.sendChan, err
	}
}

// Публикуем байтовый массив в очередь

type Rmq struct {
	Connection *amqp.Connection
	Channel    *amqp.Channel
	Failure    bool
}

func Connect() *Rmq {

	f := false

	// CONNECT
	conn, err := amqp.Dial("amqp://localhost:5672")
	if err != nil {
		log.Fatalf("Error (rmq): %v", err)
		f = true
	}
	 defer conn.Close()

	// CHANNEL
	chann, err := conn.Channel()
	if err != nil {
		log.Fatalf("Error (rmq): %v", err)
		f = true
	}

	rmq := &Rmq{
		Connection: conn,
		Channel:    chann,
		Failure:    f,
	}

	return rmq
}

func (rcl *RabbitClient) Publish(n string, b string) { r := false
	uuidNow,errUuid := uuid.NewV4()
	if errUuid != nil{
		log.Fatal(errUuid)
	}
	for {
		for {
			_, err := rcl.channel(false, r)
			if err == nil {
				break
			}
		}
		q, err := rcl.sendChan.QueueDeclare(
			n,
			true,
			false,
			false,
			false,
			amqp.Table{"x-queue-mode": "lazy"},
		)
		if err != nil {
			log.Println("--- failed to declare a queue, trying to resend ---")
			r = true
			continue
		}
		err = rcl.sendChan.Publish(
			"storage-manager",
			q.Name,
			false,
			false,
			amqp.Publishing{
				MessageId:    uuidNow.String(),
				DeliveryMode: amqp.Persistent,
				ContentType:  "text/plain",
				Body: []byte(b),
			})
		if err != nil {
			log.Println("--- failed to publish to queue, trying to resend ---")
			r = true
			continue
		}
		break
	}
}