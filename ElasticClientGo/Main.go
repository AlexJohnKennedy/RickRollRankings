package main

import "fmt"
import "os"
import k "RickRollRankings/ElasticClientGo/kafkatoelastic"
import "os/signal"
import "context"
import "sync"

const messageBufferNum = 1000;

func main() {
	bootstrapServ := os.Getenv("KAFKA_BOOTSTRAP_SERV");
	kafkaConsumeTopic := os.Getenv("TESTING_TOPIC");
	//kafkaProduceTopic := os.Getenv("");

	fmt.Printf("Test program started! Consuming from: %s", kafkaConsumeTopic);
	
	// Setup a channel to listen for shutdown signals from the OS
	signalChannel := make(chan os.Signal);
	signal.Notify(signalChannel, os.Interrupt);		// Configures this program to relay all interrupt signals to this channel
	
	// Setup a cancelation context, and a cleanup waitgroup for all our child threads to use.
	quitContext, cancelFunction := context.WithCancel(context.Background());
	var waitgroup sync.WaitGroup;

	// Launch the testing consumer with some simple config
	conf := k.KafkaConfig{
		TopicName: kafkaConsumeTopic,
		BootstrapServers: []string {bootstrapServ},
		ConsumerGroupID: "Test consumer 1",	
	}
	waitgroup.Add(1);
	go k.RunTestConsumerInGroup(quitContext, conf, &waitgroup);

	select {
	case sig := <-signalChannel:
		fmt.Printf("Received interrupt signal: %s. Calling cancel function now... \n", sig);
		cancelFunction();
		waitgroup.Wait();
		fmt.Println("Finished shutting everything down! Cya later :)");
		return;
	}
}