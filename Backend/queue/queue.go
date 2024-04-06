package queue

import "gopkg.in/Shopify/sarama.v1"

type QueueService interface {
	SendMessage(topic, message string) error
	Subscribe(topic string) (<-chan *sarama.ConsumerMessage, <-chan error)
	Close()
}
