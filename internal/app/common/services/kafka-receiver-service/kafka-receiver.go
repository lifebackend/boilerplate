package kafkareceiverservice

import (
	"strings"
	"time"

	"github.com/nvsco/wallet/internal/walleterror"
	"github.com/nvsco/wallet/pkg/scope"
)

const (
	TenSecondTime = time.Second * 10
)

type kafkaReceiverServiceInterface interface {
	HandleMessage(
		scope *scope.Scope,
		eventType string,
		data []byte,
	) *walleterror.AppMessage
}

func NewKafkaReceiver(
	kafkaBrokers string,
	topic string,
	groupID string,
	messageHandler kafkaReceiverServiceInterface,
) (*KafkaReaderService, error) {
	// nolint:exhaustivestruct
	dialer := &kafka.Dialer{
		Timeout:   TenSecondTime,
		DualStack: true,
	}

	brokers := strings.Split(kafkaBrokers, ",")

	// nolint:exhaustivestruct
	config := kafka.ReaderConfig{
		Brokers:     brokers,
		GroupID:     groupID,
		Topic:       topic,
		Dialer:      dialer,
		StartOffset: kafka.FirstOffset, // listen from the beginning of the queue
	}

	kafkaReader := kafka.NewReader(config)

	return &KafkaReaderService{
		client:         kafkaReader,
		messageHandler: messageHandler,
	}, nil
}

type KafkaReaderService struct {
	client         *kafka.Reader
	messageHandler kafkaReceiverServiceInterface
}

func (k *KafkaReaderService) handleMessage(scope *scope.Scope) *walleterror.AppMessage {
	message, err := k.client.FetchMessage(scope.Ctx)
	if err != nil {
		scope.Logger().Error(err)

		return walleterror.NewInternalServiceError(err)
	}

	eventType := ""
	isCorrectEventReceived := false
	for _, header := range message.Headers {
		if header.Key == "message_type" {
			isCorrectEventReceived = true
			eventType = string(header.Value)
			break
		}
	}

	if isCorrectEventReceived {
		msg := k.messageHandler.HandleMessage(scope, eventType, message.Value)
		if msg != nil {
			scope.Logger().Error(msg)

			// unable to process, return without commit
			return nil
		}
	}

	if !isCorrectEventReceived {
		scope.Logger().Warning("Unknown format of event, %w", message)
	}

	err = k.client.CommitMessages(scope.Ctx, message)
	if err != nil {
		scope.Logger().Warning("failed to commit messages: %w", err)

		return walleterror.NewInternalServiceError(err)
	}

	return nil
}

func (k *KafkaReaderService) Start(
	scope *scope.Scope,
) *walleterror.AppMessage {
	for {
		select {
		case <-scope.Ctx.Done():
			return nil
		default:
			msg := k.handleMessage(scope)
			if msg != nil {
				// communication failure, wait a bit
				time.Sleep(time.Second * 5) // nolint:gomnd
			}
		}
	}
}
