package kafka

import (
	"context"
	"crypto/tls"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/catcher"
)

type KafkaQueue struct {
	mechanism sasl.Mechanism
	broker    string
	groupId   string

	w *kafka.Writer
	r map[string]*kafka.Reader
}

func NewKafkaQueue(username, password, broker, groupId string) (service.ProcessQueue, error) {
	k := &KafkaQueue{
		broker:  broker,
		groupId: groupId,
		r:       map[string]*kafka.Reader{},
	}
	var err error
	k.mechanism, err = scram.Mechanism(scram.SHA256, username, password)
	if err != nil {
		return nil, err
	}
	k.w = &kafka.Writer{
		Addr: kafka.TCP(broker),
		Transport: &kafka.Transport{
			SASL: k.mechanism,
			TLS:  &tls.Config{},
		},
	}

	return k, nil
}

func (k *KafkaQueue) Write(ctx context.Context, topic string, message []byte) error {
	return k.w.WriteMessages(ctx, kafka.Message{Value: message, Topic: topic})
}

func (k *KafkaQueue) Consume(ctx context.Context, topic string, f func([]byte) error) error {
	r, err := k.getReader(topic)
	if err != nil {
		return err
	}
	for {
		message, err := r.ReadMessage(ctx)
		if err != nil {
			return err
		}

		if err := f(message.Value); err != nil {
			return err
		}
	}
}

func (k *KafkaQueue) Close() error {
	catcher := catcher.New()
	catcher.Catch(k.w.Close())
	for _, r := range k.r {
		catcher.Catch(r.Close())
	}
	return catcher.Resolve()
}

func (k *KafkaQueue) getReader(topic string) (*kafka.Reader, error) {
	if _, ok := k.r[topic]; !ok {
		k.r[topic] = kafka.NewReader(kafka.ReaderConfig{
			Brokers: []string{k.broker},
			Topic:   topic,
			GroupID: k.groupId,
			Dialer: &kafka.Dialer{
				SASLMechanism: k.mechanism,
				TLS:           &tls.Config{},
			},
		})
	}
	return k.r[topic], nil
}
