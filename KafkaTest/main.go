package main

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
	"github.com/sirupsen/logrus"
)

func main() {
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		logrus.Fatal("failed to dial leader:", err)
	}

	var msg [2000]kafka.Message

	for i := range msg {
		msg[i] = kafka.Message{Value: []byte(strconv.Itoa(i))}
	}

	conn.SetWriteDeadline(time.Now().Add(1 * time.Second))
	_, err = conn.WriteMessages(
		msg[0:]...,
	)
	if err != nil {
		logrus.Fatal("failed to write messages:", err)
	}

	if err := conn.Close(); err != nil {
		logrus.Fatal("failed to close writer:", err)
	}

	ch := make(chan int, 1)
	go consumer(ch)
	<-ch
}

func consumer(ch chan int) {
	topic := "my-topic"
	partition := 0

	conn, err := kafka.DialLeader(context.Background(), "tcp", "localhost:9092", topic, partition)
	if err != nil {
		logrus.Fatal("failed to dial leader:", err)
	}

	conn.SetReadDeadline(time.Now().Add(1 * time.Second))
	batch := conn.ReadBatch(10e3, 1e6)

	b := make([]byte, 10e4)
	pgxBatch := &pgx.Batch{}
	for {
		n, err := batch.Read(b)
		if err != nil {
			break
		}
		num := string(b[:n])
		i, _ := strconv.Atoi(num)
		pgxBatch.Queue("INSERT INTO users(number) VALUES ($1)", i)
		fmt.Println(num)
	}

	postgresPool, err := pgxpool.New(context.Background(), fmt.Sprintf("postgres://postgres:postgres@localhost:5432/postgres"))
	if err != nil {
		logrus.Fatalf("Could not connect to db %s", err)
	}
	defer postgresPool.Close()

	batchResults := postgresPool.SendBatch(context.Background(), pgxBatch)
	err = batchResults.Close()

	if err := batch.Close(); err != nil {
		logrus.Fatal("failed to close batch:", err)
	}

	if err := conn.Close(); err != nil {
		logrus.Fatal("failed to close connection:", err)
	}

	ch <- 1
}
