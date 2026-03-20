package kafka

import (
	"log"

	"github.com/IBM/sarama"
)

type Producer struct {
	asyncProducer sarama.AsyncProducer
}

func NewProducer(brokers []string) (*Producer, error) {
	producer, err := sarama.NewAsyncProducer(brokers, NewConfig())
	if err != nil {
		return nil, err
	}

	p := &Producer{asyncProducer: producer}

	// handle success
	go func() {
		for success := range producer.Successes() {
			log.Printf("Message sent to partition=%d offset=%d",
				success.Partition, success.Offset)
		}
	}()

	// handle error
	go func() {
		for err := range producer.Errors() {
			log.Println("Failed to send message:", err)
		}
	}()

	return p, nil
}

func (p *Producer) Publish(topic string, message []byte) {
	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(message),
	}
	p.asyncProducer.Input() <- msg
}

func (p *Producer) Close() error {
	return p.asyncProducer.Close()
}
