package scope

import (
	"context"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	prometheusmetrics "github.com/nvsco/wallet/tools/prometheus-metrics"
	"github.com/sirupsen/logrus"
)

func NewScope(
	rootCtx context.Context,
	logger *logrus.Entry,
	xCorrelationID string,
	metrics *prometheusmetrics.Metrics,
) *Scope {
	ctx, cancel := context.WithCancel(rootCtx)

	return &Scope{
		logger:         logger,
		xCorrelationID: xCorrelationID,
		Ctx:            ctx,
		cancel:         cancel,
		level:          "",
		resource:       "",
	}
}

func NewNoopLogger() *logrus.Entry {
	logger := logrus.New()
	logger.SetOutput(ioutil.Discard)

	return logrus.NewEntry(logger)
}

type Scope struct {
	level          string
	logger         *logrus.Entry
	xCorrelationID string
	resource       string
	Ctx            context.Context // nolint:containedctx
	cancel         context.CancelFunc
}

func (s *Scope) Finish() {
	s.cancel()
}

func (s *Scope) Logger() *logrus.Entry {
	return s.logger
}

func (s *Scope) Context() context.Context {
	return s.Ctx
}

func (s *Scope) XCorrelationID() string {
	return s.xCorrelationID
}

func (s *Scope) SetResource(resource string) {
	s.resource = resource
	s.logger = s.logger.WithField("resource", resource)
}

func (s *Scope) WithField(name string, value interface{}) *Scope {
	if name == "level" {
		return nil
	}

	return &Scope{
		level:          s.level,
		logger:         s.logger.WithField(name, value),
		xCorrelationID: s.xCorrelationID,
		Ctx:            s.Ctx,
		cancel:         s.cancel,
		resource:       s.resource,
	}
}

func (s *Scope) Child(obj interface{}) *Scope {
	level := fmt.Sprintf("%s:%s", s.level, getType(obj))

	return &Scope{
		level:          level,
		logger:         s.logger.WithField("level", level),
		xCorrelationID: s.xCorrelationID,
		Ctx:            s.Ctx,
		cancel:         s.cancel,
		resource:       s.resource,
	}
}

func (s *Scope) WithCancel() (*Scope, context.CancelFunc) {
	ctx, cancel := context.WithCancel(s.Ctx)

	return &Scope{
		level:          s.level,
		logger:         s.logger,
		xCorrelationID: s.xCorrelationID,
		Ctx:            ctx,
		cancel:         cancel,
		resource:       s.resource,
	}, cancel
}

func (s *Scope) WithTimeout(timeoutDuration time.Duration) (*Scope, context.CancelFunc) {
	ctx, cancelFunc := context.WithTimeout(s.Ctx, timeoutDuration)

	return &Scope{
		level:          s.level,
		logger:         s.logger,
		xCorrelationID: s.xCorrelationID,
		Ctx:            ctx,
		cancel:         cancelFunc,
		resource:       s.resource,
	}, cancelFunc
}

func (s *Scope) Fork() *Scope {
	ctx, cancel := context.WithCancel(context.Background())

	return &Scope{
		level:          s.level,
		logger:         s.logger,
		xCorrelationID: s.xCorrelationID,
		Ctx:            ctx,
		cancel:         cancel,
		resource:       s.resource,
	}
}

func (s *Scope) ForkWithCtx(rootCtx context.Context) *Scope {
	ctx, cancel := context.WithCancel(rootCtx)

	return &Scope{
		level:          s.level,
		logger:         s.logger,
		xCorrelationID: s.xCorrelationID,
		Ctx:            ctx,
		cancel:         cancel,
		resource:       s.resource,
	}
}

func getType(obj interface{}) string {
	name := fmt.Sprintf("%T", obj)
	parts := strings.Split(name, ".")

	return parts[len(parts)-1]
}
