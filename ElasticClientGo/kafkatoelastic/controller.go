package kafkatoelastic

import (
	"fmt"
	"sync"
	"context"
	"github.com/segmentio/kafka-go"
)

// KafkaConfig contains config information telling us how to connect to Kafka, which topic to consume from, and what partitions
// to consume from within that topic.
type KafkaConfig struct {
	topicName string
	bootstrapServers []string
	// TODO: partitionIds []int	// TODO: Make this system work on Reader-partition pairs, instead of buffering messages from all partitions into a single buffer.
}
// ElasticClientConfig contains config information tellins us how to connect to Elastic, and a function defining how to convert an
// incoming kafka message into something which can be pushed to an Elasticsearch index.
type ElasticClientConfig struct {
	elasticServers []string
	index string
	conversionFunc KafkaBytesToElasticPayloadAdaptor
}
// ElasticPayload contains information to be published to a specified elasticsearch index.
type ElasticPayload struct {
	data string
	// TODO ???
}
// KafkaBytesToElasticPayloadAdaptor is a function signature for conversion functions, which convert incoming kafka messages into a payload
// which our elasticsearch client knows how to publish.
type KafkaBytesToElasticPayloadAdaptor func(*kafka.Message) (*ElasticPayload, error);



// Run is the function you will 'run' (go-routine) to boot up the overall process of consuming messages, and
// publishing them into elastic. It can be cancelled with a context.
func Run(ctx context.Context, kafkaConfig KafkaConfig, elasticConfig ElasticClientConfig) {
	// TODO
}