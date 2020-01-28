package kafkatoelastic

import (
	"github.com/Shopify/sarama"
	//"time"
)

// ElasticClientConfig contains config information tellins us how to connect to Elastic, and a function defining how to convert an
// incoming kafka message into something which can be pushed to an Elasticsearch index.
type ElasticClientConfig struct {
	elasticServers []string
	username string
	password string
	index string
	buffersize int				// Primitive: Just the number of messages rather than an actual data-size amount.
	minimumFlushInterval int64	// ms. Defines the period at which the buffer is mandatorily flushed to elastic, and offsets committed.
	conversionFunc KafkaBytesToElasticPayloadAdaptor
}
// ElasticPayload contains information to be published to a specified elasticsearch index.
type ElasticPayload struct {
	data string
	// TODO ???
}
// KafkaBytesToElasticPayloadAdaptor is a function signature for conversion functions, which convert incoming kafka messages into a payload
// which our elasticsearch client knows how to publish.
type KafkaBytesToElasticPayloadAdaptor func(*sarama.ConsumerMessage) (*ElasticPayload, error)