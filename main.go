// consumer project main.go
package main

import (
	"encoding/json"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"os"
)

type TestData struct {
	ID    string
	Value int     `json:",string"`
	DT    float64 `json:",string"`
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
		panic(fmt.Sprintf("%s: %s", msg, err))
	}
}

func main() {
	//create and open file
	os.Mkdir("/tmp", 0777)
	err := os.Remove("tmp/stratagem_golang_output.txt")
	f, err := os.Create("/tmp/stratagem_golang_output.txt")
	failOnError(err, "Failed to create file")
	defer f.Close()

	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"test_queue", // name
		false,        // durable
		false,        // delete when usused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			myid := d.MessageId
			deltime := d.Timestamp
			const layout = "Mon Jan 2 15:04:05 -0700 MST 2006"
			strTime := deltime.UTC().Format(layout)

			var m TestData
			err := json.Unmarshal([]byte(d.Body), &m)
			failOnError(err, "Failed to unpack json")

			res := &TestData{}
			json.Unmarshal([]byte(d.Body), &res)
			//fmt.Println(res._id)
			iVal := int(res.Value)
			if iVal == 1 {
				myline := strTime + " | message ID " + myid + " | Got a 1!\n"
				line := []byte(myline)
				_, err := f.Write(line)
				failOnError(err, "Failed to write to file ")
				f.Sync()

			}

			//log.Printf("Received a message: %s", d.Body)

		}
	}()

	log.Printf(" [*] Waiting for messages New. To exit press CTRL+C")
	<-forever

}
