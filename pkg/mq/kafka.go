package kafka

import (
	"context"
	"errors"
	"log"
	"strings"

	"github.com/Shopify/sarama"
)

type WorkerFunc func(*sarama.ConsumerMessage) error

var (
	globalKafka               *Kafka
	PushMessageValueTypeError = errors.New("value type error")
)

type Config struct {
	BrokerList, Topic, GroupId string
}

type Kafka struct {
	producer sarama.SyncProducer
	consumer sarama.ConsumerGroup
	topic    string
}

func GetGlobalKafka() *Kafka {
	return globalKafka
}
func (k *Kafka) Close() {
	k.producer.Close()
	k.consumer.Close()
}
func NewKafka(c *Config, ops ...Option) (*Kafka, error) {
	dOps := defaultOptions
	for _, v := range ops {
		v(&dOps)
	}
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = dOps.producerRequiredAcks // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = dOps.producerRetryMax        // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = dOps.producerReturnSuccess
	brokerList := strings.Split(c.BrokerList, ",")
	producer, err := sarama.NewSyncProducer(brokerList, config)
	if err != nil {
		return nil, err
	}
	config.Consumer.Group.Rebalance.GroupStrategies = dOps.consumerGroupStrategies
	config.Consumer.Offsets.Initial = dOps.consumerOffsetsInit
	consumer, err := sarama.NewConsumerGroup(brokerList, c.GroupId, config)
	if err != nil {
		return nil, err
	}
	k := &Kafka{
		producer: producer,
		consumer: consumer,
		topic:    c.Topic,
	}
	globalKafka = k
	return k, nil
}
func (k *Kafka) Push(value interface{}) error {
	var messageValue sarama.Encoder
	switch value.(type) {
	case string:
		messageValue = sarama.StringEncoder(value.(string))
	case []byte:
		messageValue = sarama.ByteEncoder(value.([]byte))
	default:
		return PushMessageValueTypeError
	}
	message := &sarama.ProducerMessage{
		Topic: k.topic,
		Value: messageValue,
	}
	_, _, err := k.producer.SendMessage(message)
	return err
}
func (k *Kafka) Consumer(ctx context.Context, handler sarama.ConsumerGroupHandler) error {
	return k.consumer.Consume(ctx, []string{k.topic}, handler)
}

type DefaultConsumer struct {
	Ready chan struct{}
}

func (d *DefaultConsumer) Setup(sarama.ConsumerGroupSession) error {
	return nil
}
func (d *DefaultConsumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (d *DefaultConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for message := range claim.Messages() {
		log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		session.MarkMessage(message, "")
	}
	return nil
}
