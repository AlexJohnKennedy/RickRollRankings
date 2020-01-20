package main

import "fmt"
import "os"
import "RickRollRankings/LinkChecker/kafkaconsumer"
import "os/signal"
import "sync"

// This is the 'main' executable of this simple micro-service.
// It will import the custom 'consumer', 'condition checker', and 'producer' packages, and
// run them in three go routines each, indefinitely!

// Consumer --> Simply harvests incoming topic messages from kafka, and dumps them into a channel
// Condition checker --> recieves messages from a channel, and launches HTTP requests in child threads
//                       to look for 301/302 redirect responses which go to a known rick roll url. Any
//                       matches get pushed into an output channel.
// Producer --> Simply listens on a channel for 'confirmed rick rolls' and publishes them to kafka
func main() {
	kafkaConsumeTopic := os.Getenv("KAFKA_CHECK_REDIRECT_TOPIC_NAME");
	kafkaProduceTopic := os.Getenv("KAFKA_NEW_RICKROLL_TOPIC_NAME");
	bootstrapServ := os.Getenv("KAFKA_BOOTSTRAP_SERV");

	fmt.Printf("Producing to: %s, Consuming from: %s", kafkaProduceTopic, kafkaConsumeTopic);
	
	// Setup a channel to listen for shutdown signals from the OS
	signalChannel := make(chan os.Signal);
	signal.Notify(signalChannel, os.Interrupt);		// Configures this program to relay all interrupt signals to this channel

	// Initialise some communication channels, and then launch the consumer in a thread
	incomingLinks := make(chan string);
	quitSignal := make(chan interface{});
	var waitgroup sync.WaitGroup;
	waitgroup.Add(1);
	go kafkaconsumer.SpinUpConsumer(kafkaConsumeTopic, bootstrapServ, incomingLinks, quitSignal, &waitgroup);

	for {
		select {
		case interrupt := <-signalChannel:
			fmt.Println("Received interrupt signal! Sending shutdown signal now... ");
			quitSignal <- interrupt;
			waitgroup.Wait();
			fmt.Println("Finished shutting everything down! Cya later :)");
			return;
		case message := <-incomingLinks:
			// TODO: Pump this string into the Condition checker logic
			fmt.Println(message);
		}
	}
}
