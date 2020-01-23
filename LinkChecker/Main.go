package main

import "fmt"
import "os"
import k "RickRollRankings/LinkChecker/kafkaconsumerproducer"
import "RickRollRankings/LinkChecker/conditionchecker"
import m "RickRollRankings/LinkChecker/message"
import "os/signal"
import "context"
import "sync"

const messageBufferNum = 1000;

func main() {
	targetRickRollStrings := []string {
		// Youtube videos
		"oHg5SJYRHA0",
		"dQw4w9WgXcQ",
		"xfr64zoBTAQ",
		"dPmZqsQNzGA",
		"r8tXjJL3xcM",
		"6-HUgzYPm9g",
	
		// Fake news websites which contains an embedded Rick-roll on the site
		"latlmes.com/breaking",
		"rickrolled.com",
	};
	program(targetRickRollStrings);
}

// This is the 'main' executable of this simple micro-service.
// It will import the custom 'consumer', 'condition checker', and 'producer' packages, and
// run them in three go routines each, indefinitely!

// Consumer --> Simply harvests incoming topic messages from kafka, and dumps them into a channel
// Condition checker --> recieves messages from a channel, and launches HTTP requests in child threads
//                       to look for 301/302 redirect responses which go to a known rick roll url. Any
//                       matches get pushed into an output channel.
// Producer --> Simply listens on a channel for 'confirmed rick rolls' and publishes them to kafka
func program(targs []string) {
	kafkaConsumeTopic := os.Getenv("KAFKA_CHECK_REDIRECT_TOPIC_NAME");
	kafkaProduceTopic := os.Getenv("KAFKA_NEW_RICKROLL_TOPIC_NAME");
	bootstrapServ := os.Getenv("KAFKA_BOOTSTRAP_SERV");

	fmt.Printf("Producing to: %s, Consuming from: %s", kafkaProduceTopic, kafkaConsumeTopic);
	
	// Setup a channel to listen for shutdown signals from the OS
	signalChannel := make(chan os.Signal);
	signal.Notify(signalChannel, os.Interrupt);		// Configures this program to relay all interrupt signals to this channel
	
	// Setup a cancelation context, and a cleanup waitgroup for all our child threads to use.
	quitContext, cancelFunction := context.WithCancel(context.Background());
	var waitgroup sync.WaitGroup;

	// Initialise some communication channels, and then launch the consumer in a thread
	waitgroup.Add(1);
	incomingLinks := make(chan *m.Message, messageBufferNum);
	go k.SpinUpConsumer(quitContext, kafkaConsumeTopic, bootstrapServ, incomingLinks, &waitgroup);

	// Initialise some communication channels for the link checker, and then launch in a thread
	waitgroup.Add(1);
	input := make(chan *m.Message, 100);
	matches := make(chan *m.Message, 100);
	fails := make(chan *m.Message, 100);
	go conditionchecker.LaunchMessageChecker(quitContext, targs, input, matches, fails, &waitgroup);

	// Launch the kafka producer worker thread as well!
	go k.SpinUpProducer(quitContext, kafkaProduceTopic, bootstrapServ, matches, &waitgroup);

	for {
		select {
		case interrupt := <-signalChannel:
			fmt.Printf("Received interrupt signal: %s. Calling cancel function now... \n", interrupt);
			cancelFunction();
			waitgroup.Wait();
			fmt.Println("Finished shutting everything down! Cya later :)");
			return;
		case message := <-incomingLinks:
			fmt.Printf("Got message by user '%s' and %d links to check!\n---\nComment body: %s\n", message.AuthorName, len(message.LinksToCheck), message.CommentText);
			fmt.Printf("We will need to search the following links:\n");
			for _, s := range message.LinksToCheck {
				fmt.Println(s);
			}
			fmt.Printf("===========================================\n");

			// Now that we have done our debug prints, let's forward this message to our concurrently-running link checker
			input <- message;
		case nonmatch := <-fails:
			fmt.Printf("NO MATCH: Message by user '%s' was NOT a match :(\n", nonmatch.AuthorName);
		}
	}
}

func shittyTesting(targs []string) {
	fmt.Println("Starting");

	signalChannel := make(chan os.Signal);
	signal.Notify(signalChannel, os.Interrupt);

	var waitgroup sync.WaitGroup;
	waitgroup.Add(1);

	quitContext, cancelFunction := context.WithCancel(context.Background());

	input := make(chan *m.Message, 100);
	matches := make(chan *m.Message, 100);
	fails := make(chan *m.Message, 100);

	go conditionchecker.LaunchMessageChecker(quitContext, targs, input, matches, fails, &waitgroup);

	for {
		select {
		case interrupt := <-signalChannel:
			fmt.Printf("Received interrupt signal: %s. Calling cancel function now... \n", interrupt);
			cancelFunction();
			waitgroup.Wait();
			fmt.Println("Finished shutting everything down! Cya later :)");
			return;
		case m := <-matches:
			fmt.Printf("-------\n=============> FOUND MATCH: %s\n-------\n", m.CommentText);
		case m := <-fails:
			fmt.Printf("*******\n=============> NO MATCH: %s\n*******\n", m.CommentText);
		}
	}
}