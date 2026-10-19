package main

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

// Connect to the specified topic and partition in the server
func connect(topic string, partition int) (*kafka.Conn, error) {
	conn, err := kafka.DialLeader(
		context.Background(),
		"tcp",
		os.Getenv("KAFKA_ADDRESS"),
		topic,
		partition,
	)
	if err != nil {
		fmt.Println("[INTERNAL] Dailed to dial leader, message:", err)
	}
	return conn, err
} //end connect

// Converts a reading to byte and send's it to kafka
func writeReading(conn *kafka.Conn, reading Reading) error {
	var err error

	err = conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
	if err != nil {
		fmt.Println("[INTERNAL] Dailed to marshal message:", err)
		return errors.New(err.Error())
	}

	msg, err := json.Marshal(reading)

	if err != nil {
		fmt.Println("[INTERNAL] Dailed to marshal message:", err)
		return errors.New(err.Error())
	}

	_, err = conn.WriteMessages(
		kafka.Message{Value: []byte(msg)},
	)
	if err != nil {
		fmt.Println("[INTERNAL] Dailed to write message:", err)
		return errors.New(err.Error())
	}
	return nil
} //end writeMessages
