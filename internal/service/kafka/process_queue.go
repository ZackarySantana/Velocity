package kafka

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"os"

	"github.com/segmentio/kafka-go"
	"github.com/segmentio/kafka-go/sasl/scram"
	"github.com/zackarysantana/velocity/internal/service"
	"github.com/zackarysantana/velocity/src/catcher"
)

type processQueue struct {
	transport kafka.RoundTripper
	dialer    *kafka.Dialer
	broker    string
	groupId   string

	w *kafka.Writer
	r map[string]*kafka.Reader

	closed bool
}

type ProcessQueueConfig struct {
	Username string
	Password string
	Broker   string
	GroupId  string
}

func NewProcessQueueConfigFromEnv(groupId string) *ProcessQueueConfig {
	return &ProcessQueueConfig{
		Username: os.Getenv("KAFKA_USERNAME"),
		Password: os.Getenv("KAFKA_PASSWORD"),
		Broker:   os.Getenv("KAFKA_BROKER"),
		GroupId:  groupId,
	}
}

func NewProcessQueue(config *ProcessQueueConfig) (service.ProcessQueue, error) {
	k := &processQueue{
		broker:  config.Broker,
		groupId: config.GroupId,
		r:       map[string]*kafka.Reader{},
	}
	// If this isn't dev mode, we need to set up the SASL mechanism.
	// If it is, we can just use the default transport and dialer.
	if os.Getenv("DEV_SERVICES") != "true" {
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

func (k *processQueue) Write(ctx context.Context, topic string, messages ...[]byte) error {
	if k.closed {
		return errors.New("queue is closed")
	}

	kafkaMessages := make([]kafka.Message, len(messages))
	for i, message := range messages {
		kafkaMessages[i] = kafka.Message{Value: message, Topic: topic, Key: []byte(fmt.Sprintf("%d", i))}
	}
	return k.w.WriteMessages(ctx, kafkaMessages...)
}

func (k *processQueue) Consume(ctx context.Context, topic string, consumerFunc func(message []byte) (processed bool, err error)) error {
	if k.closed {
		return errors.New("queue is closed")
	}

	r, err := k.getReader(topic)
	if err != nil {
		return err
	}

	for {
		if k.closed {
			return errors.New("queue is closed")
		}

		message, err := r.FetchMessage(ctx)
		if err != nil {
			if errors.Is(err, io.EOF) {
				continue
			}
			return err
		}

		processed, err := consumerFunc(message.Value)
		if err != nil {
			return err
		}
		if processed {
			r.CommitMessages(ctx, message)
		}
	}
}

func (k *processQueue) Close() error {
	if k.closed {
		return errors.New("queue already closed")
	}
	catcher := catcher.New()
	catcher.Catch(k.w.Close())
	for _, r := range k.r {
		catcher.Catch(r.Close())
	}
	k.closed = true
	return catcher.Resolve()
}

func (k *processQueue) getReader(topic string) (*kafka.Reader, error) {
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
