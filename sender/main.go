package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/streadway/amqp" // BSD 2-Clause "Simplified" License: https://github.com/streadway/amqp/blob/master/LICENSE
)

func main() {

	conn, err := mqConnect("amqp://rabbitmq:5672")
	if err != nil {
		log.Fatalf("error connecting to rabbitmq: %v", err)
	}

	http.HandleFunc("/marco", func(w http.ResponseWriter, r *http.Request) {

		queue := "PoolQueue"

		ch, err := openChannel(conn, queue)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		msg := amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte("Marco"),
		}

		if err = ch.Publish("", queue, true, false, msg); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusAccepted)
		log.Println("Marco!")
		w.Write([]byte(fmt.Sprintf("Message sent on default exchange to PoolQueue at %v", time.Now())))
	})

	http.ListenAndServe(":8080", nil)
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
