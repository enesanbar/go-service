package producer

import (
	"errors"
)

var (
	ErrProducerNotFound = errors.New("producer not found")
)

// create factory of factories for RabbitMQ, AWS SQS, and Kafka producers
// type ProducerFactory interface {
// 	CreateProducer(name string) (Producer, error)
// }
