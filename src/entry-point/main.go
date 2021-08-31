package main

import (
	"zeithub.com/site-status-controller/api/gRPC/rmq"
)

func main() {

	var rc rmq.RabbitClient
	rmq.Connect()
	rc.Publish("Orkhan lox","test-queue")

	//rc.Consume("url_array")
}

