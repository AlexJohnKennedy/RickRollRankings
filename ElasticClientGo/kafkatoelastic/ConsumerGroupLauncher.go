package kafkatoelastic

import (
	"fmt"
	"context"
	"sync"
	"github.com/Shopify/sarama"
)

// KafkaConfig contains config information telling us how to connect to Kafka, which topic to consume from, and what partitions to consume from within that topic.
type KafkaConfig struct {
	TopicName string
	BootstrapServers []string
	ConsumerGroupID string
}

// RunTestConsumerInGroup memes
func RunTestConsumerInGroup(ctx context.Context, conf KafkaConfig, quitWg *sync.WaitGroup) {
	defer func() {
		fmt.Println("ConsumerGroup launcher thread shutting down now...");
		quitWg.Done();
	}();

	// Init config, specify appropriate version
	config := sarama.NewConfig()
	config.Version = sarama.V2_1_0_0

	// Start a client: Connects to the cluster via the specified bootstrap servers
	client, err := sarama.NewClient(conf.BootstrapServers, config);
	if err != nil {
		panic(err);
	}
	defer func() { _ = client.Close(); }();

	// Start a new consumer group: Uses the client to connect to a consumer group
	consumerGroup, err := sarama.NewConsumerGroupFromClient(conf.ConsumerGroupID, client);
	if err != nil {
		panic(err);
	}
	defer func() { _ = consumerGroup.Close(); }();

	// Track errors. Sarama outputs consumer-group errors to a specific channel. We will simply print them
	go func() {
		for err := range consumerGroup.Errors() {
			fmt.Println("ERROR ", err);
		}
	}();

	// Instantiate a new handler for the consumer group. We only need a single handler for the group, but that handler will have it's
	// "consume" hook invoked in parallel on multiple goroutines, one for each "claim" assigned by kafka. The number of assignments depends
	// on how many other clients are connecting to kafka under the same consumer group id! (and of the number of partitions in the first place).
	handler := ClaimHandler{
		consumerClientName: "CLIENT A",
		mutex: &sync.Mutex{},
		sessionCounter: 0,
	}

	// Iterate over consumer sessions. I.e., each time this is called, a new consumergroupsession is kicked off, and this thread is blocked
	// for the duration of that session. When kafka triggers a rebalance operation then the session is closed, and it is up to us to open 
	// another one. Hence, we are triggering these consumer sessions in a loop until the context is cancelled.
	for {
		err := consumerGroup.Consume(ctx, []string {conf.TopicName}, &handler);
		if err != nil {
			fmt.Println("MEGA ERROR DURING ConsumerGroup.Consumer(): ", err);
		}
		if ctx.Err() != nil {
			fmt.Println("Context cancelled! Shutting down now... ", err);
			return;
		}
	}
}