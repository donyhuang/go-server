package kafka

import (
	"github.com/Shopify/sarama"
	"sync"
	"time"
)

type BatchWorkerFunc func([]*sarama.ConsumerMessage) error
type BatchConsumer struct {
	tick   *time.Ticker
	data   []*sarama.ConsumerMessage
	worker BatchWorkerFunc
	sync.Mutex
	session sarama.ConsumerGroupSession
}

func (b *BatchConsumer) Setup(session sarama.ConsumerGroupSession) error {
	return nil
}

func (b *BatchConsumer) Cleanup(session sarama.ConsumerGroupSession) error {
	b.tick.Stop()
	return nil
}

func (b *BatchConsumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	b.session = session
	for {
		select {
		case message := <-claim.Messages():
			b.Lock()
			b.data = append(b.data, message)
			b.Unlock()
		case <-session.Context().Done():
			return nil
		}
	}
}
