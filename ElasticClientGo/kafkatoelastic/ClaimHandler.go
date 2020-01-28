package kafkatoelastic

import "fmt"
import "strings"
import "github.com/Shopify/sarama"
import "sync"

// Define a consumer-claim-handler custom type.
// This acts as the logic which processes incoming messages from a consumer group, and then decides
// how to commit/rollback those offsets afterwards.
// We provide callbacks for lifecycle functions which is invoked by the sarama lifecycle.
// In the context of sarama, a "claim" is an active assignment to particular partition of a topic to be consumed from.
// The ConsumerGroup works by recieving an assignment of one-or-more claims, which is then "consumed" by the ConsumeClaim
// function. This function is presumably invoked once for each claim inside the "session", where the "session" represents this
// particular session-of-time under a particular partition-assignment for this consumer client. I.e., when the kafka-cluster
// decides to re-assign partitions for our consumergroup (say, another consumer in the consumer group died) then the session
// ends and another one is started once the rebalancing is completed. This is what the 'setup'/'cleanup' lifecycle methods are
// for!
// Remember, one instance of a "consumer group" in the sarama library in the context of this program is actually a single
// "consumer host", in the specified consumer group. Of which there are possibly many others, runnning as completely different
// programs on other machines! Thus, when a re-balance happens, this program will recieve new "claims" which informs the sarama
// library which data to request out of kafka during the session.

// ----------------------------------------------------------------------------------------------------------------------------

// ClaimHandler is our custom claim handler.
type ClaimHandler struct {
	consumerClientName string	// Passed in as an environment var to test the consumer-group assignment behaviour when running this program multiple times as different processes.
	mutex *sync.Mutex
	sessionCounter int			// Each 'consume claim' call will grab this and increment it atomically, for debugging purposes to print what is going on!
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (h *ClaimHandler) Setup(session sarama.ConsumerGroupSession) error {
	// Print the setup information in one go, using a string builder and Fprintf
	var b strings.Builder;
	fmt.Fprintf(&b, "===> NEW SESSION: %s (Generation %2d, Cluster member = '%s') <===\n", h.consumerClientName, session.GenerationID(), session.MemberID());
	fmt.Fprintf(&b, "Setting up with the following claims:\n");
	i := 1;
	for topic, partitions := range session.Claims() {
		for partition := range partitions {
			fmt.Fprintf(&b, "   %2d - Topic = '%s' Partition = %2d\n", i, topic, partition);
			i++;
		}
	}
	for j := 0; j < len(h.consumerClientName) + len(session.MemberID()); j++ {
		b.WriteRune('=');
	}
	fmt.Fprintf(&b, "============================================================\n");
	fmt.Println(b.String());
	return nil;
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (h *ClaimHandler) Cleanup(session sarama.ConsumerGroupSession) error {
	fmt.Printf("**** END SESSION: %s (Generation %2d, Cluster member = '%s') ****\n", h.consumerClientName, session.GenerationID(), session.MemberID());

	h.mutex.Lock();
	defer h.mutex.Unlock();
	h.sessionCounter = 0;
	
	return nil;
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (h *ClaimHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// Acquire a counter value for this 'consumption'
	h.mutex.Lock();
	counter := h.sessionCounter;
	h.sessionCounter++;
	h.mutex.Unlock();

	fmt.Printf("    ===> ConsumeClaim started! ID = %2d, Topic = '%s', Partition = %2d <===\n", counter, claim.Topic(), claim.Partition());
	
	// Range over the message channel. This will produce individual messages from this particular partition. When kafka decides to do a rebalance,
	// the message channel will be closed and we must clean up after ourselves before the 'rebalance timeout' time expires.
	for msg := range claim.Messages() {
		var b strings.Builder;
		fmt.Fprintf(&b, "     ---> %2d received message: Topic = '%s', Partitions = %2d, Offset = %6d \n", counter, msg.Topic, msg.Partition, msg.Offset);
		b.WriteString("    ");
		b.Write(msg.Value);
		fmt.Println(b.String());
	}

	// Channel closed! This means our session is about to end!
	fmt.Printf("    **** ConsumeClaim ended! ID = %d ****\n", counter);

	return nil;
}