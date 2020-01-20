package kafkaconsumer

import (
	"fmt"
	//"github.com/segmentio/kafka-go"
)

// SpinUpConsumer is a simple function returning a string-channel. The purpose of this function is to 
// pull data from a kafka topic, interpret it as a string, and output it to a channel.
// Internally, the function call spins up a kafka consumer and listens to it until a quit signal is 
// recieved (via a 'quit' channel)
func SpinUpConsumer(topicName string, output chan string, quitSignal chan interface {}) {
	fmt.Printf("Test")
}