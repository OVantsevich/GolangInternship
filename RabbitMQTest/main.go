package main

import (
	"context"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"strconv"
)

func main() {
	conn, err := amqp.Dial("amqp://rabbit:rabbit@localhost:5672/")
	if err != nil {
		logrus.Fatal("failed to dial:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal("failed to open a channel:", err)
	}

	q, err := ch.QueueDeclare(
		"my-topic", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		logrus.Fatal("failed to declare a queue:", err)
	}
	defer ch.Close()

	for i := 0; i < 2000; i++ {
		err = ch.PublishWithContext(
			context.Background(), // exchange
			"",
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(strconv.Itoa(i)),
			})
		if err != nil {
			logrus.Fatal("failed to publish:", err)
		}
	}

	c := make(chan struct{}, 1)
	go consumer(c)
	<-c
}

func consumer(c chan struct{}) {
	conn, err := amqp.Dial("amqp://rabbit:rabbit@localhost:5672/")
	if err != nil {
		logrus.Fatal("failed to dial:", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		logrus.Fatal("failed to open a channel:", err)
	}

	q, err := ch.QueueDeclare(
		"my-topic", // name
		false,      // durable
		false,      // delete when unused
		false,      // exclusive
		false,      // no-wait
		nil,        // arguments
	)
	if err != nil {
		logrus.Fatal("failed to declare a queue:", err)
	}
	defer ch.Close()

	ch.Qos(2000, 0, true)
	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	if err != nil {
		logrus.Fatal("failed to register a consumer", err)
	}

	var forever chan struct{}

	go func() {
		for d := range msgs {
			logrus.Printf("Received a message: %s", d.Body)
		}
	}()
	<-forever

	//b := make([]byte, 10e4)
	//pgxBatch := &pgx.Batch{}
	//for {
	//	n, err := batch.Read(b)
	//	if err != nil {
	//		break
	//	}
	//	num := string(b[:n])
	//	i, _ := strconv.Atoi(num)
	//	pgxBatch.Queue("INSERT INTO users(number) VALUES ($1)", i)
	//	fmt.Println(num)
	//}
	//
	//postgresPool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://postgres:postgres@localhost:5432/postgres"))
	//if err != nil {
	//	logrus.Fatalf("Could not connect to db %s", err)
	//}
	//defer postgresPool.Close()
	//
	//batchResults := postgresPool.SendBatch(context.Background(), pgxBatch)
	//err = batchResults.Close()
	//
	//if err := batch.Close(); err != nil {
	//	logrus.Fatal("failed to close batch:", err)
	//}
	//
	//if err := conn.Close(); err != nil {
	//	logrus.Fatal("failed to close connection:", err)
	//}
	//
	//c <- struct{}{}
}
