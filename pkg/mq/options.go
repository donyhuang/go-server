package kafka

import "github.com/Shopify/sarama"

var defaultOptions = options{
	producerRequiredAcks:    sarama.WaitForAll,
	producerRetryMax:        10,
	producerReturnSuccess:   true,
	consumerOffsetsInit:     sarama.OffsetOldest,
	consumerGroupStrategies: []sarama.BalanceStrategy{sarama.BalanceStrategyRoundRobin},
}

type options struct {
	producerRequiredAcks    sarama.RequiredAcks
	producerRetryMax        int
	producerReturnSuccess   bool
	consumerOffsetsInit     int64
	consumerGroupStrategies []sarama.BalanceStrategy
}

type Option func(*options)

func WithProducerAcks(acks sarama.RequiredAcks) Option {
	return func(o *options) {
		o.producerRequiredAcks = acks
	}
}
func WithProducerRetryMax(retryMax int) Option {
	return func(o *options) {
		o.producerRetryMax = retryMax
	}
}
func WithProducerReturnSuccess(returnSuccess bool) Option {
	return func(o *options) {
		o.producerReturnSuccess = returnSuccess
	}
}

func WithConsumerOffsetInit(offsets int64) Option {
	return func(o *options) {
		o.consumerOffsetsInit = offsets
	}
}
func WithConsumerBalanceStrategy(strategy []sarama.BalanceStrategy) Option {
	return func(o *options) {
		o.consumerGroupStrategies = strategy
	}
}
