package rmq

import (
	"fmt"
	uuid "github.com/nu7hatch/gouuid"
	"github.com/streadway/amqp"
	"log"
	"time"
	localConfig "zeithub.com/site-status-controller/api/config"
)

type RabbitClient struct {
	sendConn *amqp.Connection
	recConn  *amqp.Connection
	sendChan *amqp.Channel
	recChan  *amqp.Channel
}


func (rcl *RabbitClient) connect(isRec, reconnect bool) (*amqp.Connection, error) {
	config := localConfig.GetConfig()

	fmt.Sprintf("amqp://%s:%s/", config.Host)

	if reconnect {
		if isRec {
			rcl.recConn = nil
		} else {
			rcl.sendConn = nil
		}
	}
	if isRec && rcl.recConn != nil {
		return rcl.recConn, nil
	} else if !isRec && rcl.sendConn != nil {
		return rcl.sendConn, nil
	}
	var c string
	if config.User == "" {
		c = fmt.Sprintf("amqp://%s:%s/", config.Host, config.Port)
	} else {
		c = fmt.Sprintf("amqp://%s:%s@%s:%s/", config.User, config.Password, config.Host, config.Port)
	}
	conn, err := amqp.Dial(c)
	if err != nil {
		log.Printf("\r\n--- could not create a conection ---\r\n")
		time.Sleep(1 * time.Second)
		return nil, err
	}
	if isRec {
		rcl.recConn = conn
		return rcl.recConn, nil
	} else {
		rcl.sendConn = conn
		return rcl.sendConn, nil
	}
}

// Подписка


// Подписываемся исходя из имени очереди

func (rcl *RabbitClient) Consume(n string) {

	uuidNow, err4 := uuid.NewV4()
	if err4 != nil{
		log.Fatal(err4)
	}
	for {
		for {
			_, err := rcl.channel(true, true)
			if err == nil {
				break
			}
		}
		log.Printf("--- connected to consume '%s' ---\r\n", n)
		q, queueDeclareError := rcl.recChan.QueueDeclare(
			n,
			true,
			false,
			false,
			false,
			amqp.Table{"x-queue-mode": "lazy"},
		)
		if queueDeclareError != nil {
			log.Println("--- failed to declare a queue, trying to reconnect ---")
			continue
		}
		connClose := rcl.recConn.NotifyClose(make(chan *amqp.Error))
		connBlocked := rcl.recConn.NotifyBlocked(make(chan amqp.Blocking))
		chClose := rcl.recChan.NotifyClose(make(chan *amqp.Error))
		var m, err = rcl.recChan.Consume(
			q.Name,
			uuidNow.String(),
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			log.Println("--- failed to consume from queue, trying again ---")
			continue
		}
		shouldBreak := false
		for {
			if shouldBreak {
				break
			}
			select {
			case _ = <-connBlocked:
				log.Println("--- connection blocked ---")
				shouldBreak = true
				break
			case err = <-connClose:
				log.Println("--- connection closed ---")
				shouldBreak = true
				break
			case err = <-chClose:
				log.Println("--- channel closed ---")
				shouldBreak = true
				break
			case d := <-m:
				err := d.Body
				if err != nil {
					_ = d.Ack(false)
					break
				}
				_ = d.Ack(true)
			}
		}
	}
}

