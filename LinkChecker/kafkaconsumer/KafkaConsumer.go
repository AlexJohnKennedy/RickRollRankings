package kafkaconsumer

import (
	"fmt"
	"sync"
	"context"
	"github.com/segmentio/kafka-go"
)

// SpinUpConsumer is a simple function returning a string-channel. The purpose of this function is to 
// pull data from a kafka topic, interpret it as a string, and output it to a channel.
// Internally, the function call spins up a kafka consumer and listens to it until a quit signal is 
// recieved (via a 'quit' channel)
func SpinUpConsumer(ctx context.Context, topicName string, server string, output chan string, quitWg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Consumer worker is releasing wait-group...");
		quitWg.Done();
	}();
		
	// Using the kafka-go library we can setup up a 'reader' object which acts as a consumer group under the hood.
	fmt.Printf("Setting up consumer group to listen to topic: '%s'... \n", topicName);
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{ server },
		GroupID:   "LinkChecker",
		Topic:     topicName,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	});
	fmt.Println("Consumer started!");

	// Defer a clean up function which we will use to shutdown our consumer gracefully, and signal our completion, when the time comes.
	defer func() {
		fmt.Println("Shutting down consumer... ");
		reader.Close();
	}();

	// Run the consumer continuously. We will rely on the context cancellation to break us out of this blocking loop.
	for {
		message, err := reader.ReadMessage(ctx);
		if err != nil {
			fmt.Printf("ERROR ON READ: %s\n", err.Error());
			return;
		}
		output <- string(message.Value[:]);
	}
}