package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"io"
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

type KafkaQueueConfig struct {
	Username string
	Password string
	Broker   string
	GroupId  string
}

func NewKafkaQueueOptionsFromEnv() *KafkaQueueConfig {
	return &KafkaQueueConfig{
		Username: os.Getenv("KAFKA_USERNAME"),
		Password: os.Getenv("KAFKA_PASSWORD"),
		Broker:   os.Getenv("KAFKA_BROKER"),
		GroupId:  os.Getenv("KAFKA_GROUP_ID"),
	}
}

func NewKafkaQueue(config *KafkaQueueConfig) (service.ProcessQueue, error) {
	k := &KafkaQueue{
		broker:  config.Broker,
		groupId: config.GroupId,
		r:       map[string]*kafka.Reader{},
	}
	// If this isn't dev mode, we need to set up the SASL mechanism.
	// If it is, we can just use the default transport and dialer.
	if os.Getenv("DEV_MODE") != "true" {
		mechanism, err := scram.Mechanism(scram.SHA256, config.Username, config.Password)
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
		Addr:      kafka.TCP(config.Broker),
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

func (k *KafkaQueue) Consume(ctx context.Context, topic string, f func([]byte) (bool, error)) error {
	r, err := k.getReader(topic)
	if err != nil {
		return err
	}
	for {
		message, err := r.ReadMessage(ctx)
		message, err = r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			return err
		}

		ok, err := f(message.Value)
		if err != nil {
			return err
		}
		if ok {
			r.CommitMessages(ctx, message)
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
