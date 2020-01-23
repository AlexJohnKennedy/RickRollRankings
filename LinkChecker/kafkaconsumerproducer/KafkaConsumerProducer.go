package kafkaconsumerproducer

import (
	"fmt"
	"sync"
	"context"
	"github.com/segmentio/kafka-go"
	m "RickRollRankings/LinkChecker/message"
)

// SpinUpConsumer is a simple function which pulls data from a kafka topic and outputs it to a channel.
// Internally, the function call spins up a kafka consumer and listens to it until a quit signal is 
// recieved (via a 'quit' channel)
func SpinUpConsumer(ctx context.Context, topicName string, server string, output chan *m.Message, quitWg *sync.WaitGroup) {
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
		kafkaMsg, err := reader.ReadMessage(ctx);
		if err != nil {
			fmt.Printf("ERROR ON READ: %s\n", err.Error());
			return;
		}
		decodedMsg, err := m.BuildMessage(kafkaMsg.Value, kafkaMsg.Key);
		if (err == nil) {
			output <- decodedMsg;
		} else {
			fmt.Printf("ERROR ON DATA PARSE: %s\n", err.Error());
		}
	}
}

// SpinUpProducer listens on a channel and uses a kafka writer to publish those messages into a specified kafka topic.
func SpinUpProducer(ctx context.Context, topicName string, server string, input chan *m.Message, quitWg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Producer worker is releasing wait-group...");
		quitWg.Done();
	}();

	fmt.Printf("Setting up producer to publish to topic: '%s'... \n", topicName);
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: 	[]string{ server },
		Topic:	 	topicName,
		RequiredAcks: 1,
	});

	// Defer a clean up function which we will use to shutdown our consumer gracefully, and signal our completion, when the time comes.
	defer func() {
		fmt.Println("Shutting down producer... ");
		writer.Close();
	}();
	
	var producerWg sync.WaitGroup;

	for {
		select {
		case <-ctx.Done():
			producerWg.Wait();
			return;
		case message := <-input:
			fmt.Printf("Found a match!! Message by user '%s' was a match :) What a troll\n", message.AuthorName);
			if data, err := message.ToJSONWithoutLinks(); err == nil {
				producerWg.Add(1);
				go func(message *m.Message) {
					defer func() {
						producerWg.Done();
					}();
					writer.WriteMessages(ctx, kafka.Message{
						Key: message.KafkaKey,
						Value: data,
					});
				}(message);
			} else {
				fmt.Printf("Error encoding message into Json bytes: %s\n", err.Error());
			}
		}
	}
}