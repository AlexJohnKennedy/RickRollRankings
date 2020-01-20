package kafkaconsumer

import (
	"fmt"
	"sync"
	//"github.com/segmentio/kafka-go"
)

// SpinUpConsumer is a simple function returning a string-channel. The purpose of this function is to 
// pull data from a kafka topic, interpret it as a string, and output it to a channel.
// Internally, the function call spins up a kafka consumer and listens to it until a quit signal is 
// recieved (via a 'quit' channel)
func SpinUpConsumer(topicName string, server string, output chan string, quitSignal chan interface {}, quitWg *sync.WaitGroup) {
	defer func() {
		fmt.Println("Consumer worker is releasing wait-group...");
		quitWg.Done();
	}();
	
	fmt.Println("Setting up consumer...");

	// TODO: Consumer setup code

	fmt.Println("Consumer started!");

	// Defer a clean up function which we will use to shutdown our consumer gracefully, and signal our completion, when the time comes.
	defer func() {
		fmt.Println("Shutting down consumer... ");
		// TODO: Consumer shutdown code
	}();

	// TODO: Run the consumer, pumping data into the output channel. Splice in a 'select' on the quitSignal channel so we can return
	// when required.
	for {
		select {
		case <-quitSignal:
			return;
		}
	}
}