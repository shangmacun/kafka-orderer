package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"

	cluster "github.com/bsm/sarama-cluster"
)

func launchConsumer(config *cluster.Config, userPrefs *prefs) {
	consumer := newConsumer(config, userPrefs)
	defer func() {
		if err := consumer.Close(); err != nil {
			log.Fatalln(err)
		}
	}()

	// Trap SIGINT to trigger a shutdown
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	for {
		select {
		case <-signals:
			return
		case err := <-consumer.Errors():
			log.Printf("Error: %s\n", err.Error())
		case note := <-consumer.Notifications():
			log.Printf("Rebalanced: %+v\n", note)
		case message := <-consumer.Messages():
			fmt.Fprintf(os.Stdout, "Consumed message (topic: %s, part: %d, offset: %d, value: %s)\n",
				message.Topic, message.Partition, message.Offset, message.Value)
			consumer.MarkOffset(message, "")
		}
	}
}

func newConsumer(config *cluster.Config, userPrefs *prefs) *cluster.Consumer {
	brokers := strings.Split(userPrefs.brokers, ",")
	topics := []string{userPrefs.topic}
	config.Consumer.Offsets.Initial = userPrefs.begin
	consumer, err := cluster.NewConsumer(brokers, userPrefs.group, topics, config)
	if err != nil {
		panic(err)
	}
	return consumer
}