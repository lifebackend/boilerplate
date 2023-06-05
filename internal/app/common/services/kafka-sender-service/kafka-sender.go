package kafkasenderservice

import (
	"database/sql"
	"strings"
	"time"

	"github.com/nvsco/wallet/internal/app/wallet/config"

	"github.com/nvsco/wallet/internal/walleterror"
	"github.com/nvsco/wallet/pkg/scope"
)

const (
	KafkaWorkersNumber = 2
	TimeBaseFormat     = 10
	SixtySecondTime    = time.Second * 60
	TenSecondTime      = time.Second * 10
	FiveSecondTime     = time.Second * 5
	SearchLimit        = 100
)

func NewKafkaSender(cfg *config.Config) (*KafkaSenderService, error) {
	// nolint:exhaustivestruct
	dialer := &kafka.Dialer{
		Timeout:   TenSecondTime,
		DualStack: true,
	}

	kafkaWriter := kafka.NewWriter(kafka.WriterConfig{
		Brokers:           strings.Split(cfg.KafkaBrokers, ","),
		Topic:             "",
		Dialer:            dialer,
		Balancer:          nil,
		MaxAttempts:       0,
		QueueCapacity:     0,
		BatchSize:         0,
		BatchBytes:        0,
		BatchTimeout:      0,
		ReadTimeout:       0,
		WriteTimeout:      FiveSecondTime,
		RebalanceInterval: 0,
		IdleConnTimeout:   0,
		RequiredAcks:      -1,
		Async:             false,
		CompressionCodec:  nil,
		Logger:            nil,
		ErrorLogger:       nil,
	})

	return &KafkaSenderService{
		client:               kafkaWriter,
		checkImmediatelyChan: make(chan struct{}, KafkaWorkersNumber),
	}, nil
}

type KafkaSenderService struct {
	client               *kafka.Writer
	checkImmediatelyChan chan struct{}
	db                   *sql.DB
}

//func (k *KafkaSenderService) QueueEvent(
//	scope *scope.Scope,
//	topic string,
//	partitionShardingKey string,
//	version string,
//	application string,
//	eventType string,
//	createdAt time.Time,
//	messageID uuid.UUID,
//	event interface{},
//) *walleterror.AppMessage {
//	data, err := json.Marshal(event)
//	if err != nil {
//		return walleterror.NewInternalServiceError(err)
//	}
//
//	storedKafkaEvent := &dbmodels.KafkaEventsQueue{
//		ID:                   0,
//		Topic:                topic,
//		PartitionShardingKey: partitionShardingKey,
//		HeaderVersion:        version,
//		HeaderApplication:    application,
//		HeaderType:           eventType,
//		HeaderTimestamp:      strconv.FormatInt(createdAt.Unix(), TimeBaseFormat),
//		HeaderMessageID:      messageID,
//		Payload:              data,
//		CreatedAt:            time.Time{},
//		HeaderCorrelationID:  scope.XCorrelationID(),
//	}
//
//	msg := storedKafkaEvent.Insert(scope)
//	if msg != nil {
//		return msg
//	}
//
//	return nil
//}

func (k *KafkaSenderService) Start(
	scope *scope.Scope,
) *walleterror.AppMessage {
	// spawn few workers
	for i := 0; i < KafkaWorkersNumber; i++ {
		go k.processInternalKafkaEventAsync(scope)
	}

	// init the ticker to periodically re-check the queue
	t := time.NewTicker(TenSecondTime)
	defer t.Stop()

	// plan the next check immediately
	k.triggerNextEventChecking()

	for {
		select {
		case <-scope.Ctx.Done():
			return nil
		case <-t.C:
			// periodically re-check the queue
			k.triggerNextEventChecking()
		}
	}
}

func (k *KafkaSenderService) processInternalKafkaEventAsync(
	rootScope *scope.Scope,
) {
	for {
		select {
		case <-rootScope.Ctx.Done():
			return
		case <-k.checkImmediatelyChan:
			scope := rootScope.Fork()

			paxErr := k.processInternalKafkaEvent(scope)
			if paxErr != nil {
				scope.Logger().Error(paxErr)
			}

			scope.Finish()
		}
	}
}

func (k *KafkaSenderService) triggerNextEventChecking() {
	if len(k.checkImmediatelyChan) == 0 {
		select {
		case k.checkImmediatelyChan <- struct{}{}:
		default:
		}
	}
}

func (k *KafkaSenderService) processInternalKafkaEvent(
	rootScope *scope.Scope,
) *walleterror.AppMessage {
	scope, cancel := rootScope.WithTimeout(FiveSecondTime)
	defer cancel()

	tx, err := k.db.BeginTx(scope.Ctx, nil)
	if err != nil {
		rootScope.Logger().Warning(err)
		return walleterror.NewInternalServiceError(err)
	}
	defer tx.Rollback()

	paxErr := k.processKafkaEvents(scope, tx)
	if paxErr != nil {
		return paxErr
	}

	err = tx.Commit()
	if err != nil {
		return walleterror.NewInternalServiceError(err)
	}

	return nil
}

//func (k *KafkaSenderService) processKafkaEvents(
//	scope *scope.Scope,
//) *walleterror.AppMessage {
//	storedKafkaEvents, err := k.xo.GetKafkaEventsQueueWithLimitForUpdate(scope, tx, SearchLimit)
//	if err != nil {
//		return walleterror.NewInternalServiceError(err)
//	}
//
//	if len(storedKafkaEvents) == 0 {
//		return nil
//	}
//
//	if len(storedKafkaEvents) > 1 {
//		// as we have more than one task - backward notify that new task should be checked
//		k.triggerNextEventChecking()
//	}
//
//	// send events
//	paxErr := k.send(scope, storedKafkaEvents)
//	if paxErr != nil {
//		return paxErr
//	}
//
//	ids := make([]int64, 0)
//	for _, storedKafkaEvent := range storedKafkaEvents {
//		ids = append(ids, storedKafkaEvent.ID)
//	}
//
//	_, err = k.xo.DeleteKafkaEventsQueueByIds(scope, tx, dbmodels.NewInt64Slice(ids))
//	if err != nil {
//		return walleterror.NewUnknownError(err)
//	}
//
//	return nil
//}

func (k *KafkaSenderService) send(
	scope *scope.Scope,
	kafkaEvents []*dbmodels.KafkaEventsQueue,
) *walleterror.AppMessage {
	msgs := make([]kafka.Message, 0, len(kafkaEvents))

	for _, kafkaEvent := range kafkaEvents {
		msg := kafka.Message{
			Topic:     kafkaEvent.Topic,
			Partition: 0,
			Offset:    0,
			Key:       []byte(kafkaEvent.PartitionShardingKey),
			Value:     kafkaEvent.Payload,
			Headers: []kafka.Header{
				{Key: "version", Value: []byte(kafkaEvent.HeaderVersion)},
				{Key: "application_name", Value: []byte(kafkaEvent.HeaderApplication)},
				{Key: "message_type", Value: []byte(kafkaEvent.HeaderType)},
				{Key: "timestamp", Value: []byte(kafkaEvent.HeaderTimestamp)},
				{Key: "message_id", Value: []byte(kafkaEvent.HeaderMessageID.String())},
				{Key: "correlation_id", Value: []byte(kafkaEvent.HeaderCorrelationID)},
			},
			Time: time.Now(),
		}

		msgs = append(msgs, msg)
	}

	timer := time.NewTimer(SixtySecondTime)
	defer timer.Stop()

	var lastErr error
	for {
		lastErr = k.client.WriteMessages(scope.Ctx, msgs...)
		if lastErr == nil {
			break
		}

		scope.Logger().Errorf("Kafka event sending failed: %s", lastErr)

		select {
		case <-scope.Ctx.Done():
			return walleterror.NewInternalServiceError(lastErr)
		case <-timer.C:
			return walleterror.NewInternalServiceError(lastErr)
		case <-time.After(time.Second):
		}
	}

	return nil
}
