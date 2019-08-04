package main

import (
	."app/MQdemo/api"
	"app/MQdemo/model"
	"encoding/json"
	"github.com/streadway/amqp"
	"log"
	"os"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	CheckError(err, "Can't connect to RabbitMQ")
	defer conn.Close()

	channel, err := conn.Channel()
	CheckError(err, "Can't create a channel")
	defer channel.Close()

	queue, err := channel.QueueDeclare(
		"order",
		true,
		false,
		false,
		false,
		nil)
	CheckError(err, "Can't declare a queue")

	err = channel.Qos(1, 0, false)
	CheckError(err, "Can't configure the QoS")

	messageChannel, err := channel.Consume(
		queue.Name,
		"",
		false,
		false,
		false,
		false,
		nil,
	)
	CheckError(err, "Can't register a consumer")

	stopChan := make(chan bool)
	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range messageChannel {
			log.Printf("Received a message: %s", string(d.Body))
			var order model.Order
			err = json.Unmarshal(d.Body, &order)
			CheckError(err,"Error encoding JSON")

			order.AddOrder()
			order.ChangeInventory()

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %s", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	log.Println("Waiting for messages.")
	<-stopChan
}