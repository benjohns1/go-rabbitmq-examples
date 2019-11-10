package main

import (
	"log"
	"time"

	"github.com/streadway/amqp" // BSD 2-Clause "Simplified" License: https://github.com/streadway/amqp/blob/master/LICENSE
)

func main() {

	conn, err := mqConnect("amqp://rabbitmq:5672")
	if err != nil {
		log.Fatalf("error connecting to rabbitmq: %v", err)
	}

	queue := "PoolQueue"

	ch, err := openChannel(conn, queue)
	if err != nil {
		log.Fatalf("error opening channel: %v", err)
	}

	replies, err := ch.Consume(queue, "receiver", true, true, false, false, nil)
	if err != nil {
		log.Fatalf("error listening on %s", queue)
	}

	for {
		select {
		case reply := <-replies:
			if string(reply.Body) == "Marco" {
				log.Println("Polo!")
			} else {
				log.Printf("unknown message: %v\n", string(reply.Body))
			}
		}
	}
}

func mqConnect(addr string) (conn *amqp.Connection, err error) {

	const maxRetries = 60
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(addr)
		if err == nil {
			log.Printf("connected to RabbitMQ on %s", addr)
			return conn, nil
		}
		time.Sleep(3 * time.Second)
		log.Printf("retrying RabbitMQ connection after error: %v", err)
	}
	return
}

func openChannel(conn *amqp.Connection, queueName string) (*amqp.Channel, error) {

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if _, err = ch.QueueDeclare(queueName, true, true, false, false, nil); err != nil {
		return nil, err
	}

	return ch, nil
}
