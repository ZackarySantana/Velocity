package kafka

import (
	"context"
	"crypto/tls"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/catcher"
)

type KafkaQueue struct {
	transport kafka.RoundTripper
	dialer    *kafka.Dialer
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
	// If this isn't dev mode, we need to set up the SASL mechanism.
	// If it is, we can just use the default transport and dialer.
	if os.Getenv("DEV_MODE") != "true" {
		mechanism, err := scram.Mechanism(scram.SHA256, username, password)
		if err != nil {
			return nil, err
		}
		k.transport = &kafka.Transport{
			SASL: mechanism,
			TLS:  &tls.Config{},
		}
		k.dialer = &kafka.Dialer{
			SASLMechanism: mechanism,
			TLS:           &tls.Config{},
		}
	}
	k.w = &kafka.Writer{
		Addr:      kafka.TCP(broker),
		Transport: k.transport,
	}

	return k, nil
}

func (k *KafkaQueue) Write(ctx context.Context, topic string, messages ...[]byte) error {
	kafkaMessages := make([]kafka.Message, len(messages))
	for i, message := range messages {
		kafkaMessages[i] = kafka.Message{Value: message, Topic: topic}
	}
	return k.w.WriteMessages(ctx, kafkaMessages...)
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
			Dialer:  k.dialer,
		})
	}
	return k.r[topic], nil
}
