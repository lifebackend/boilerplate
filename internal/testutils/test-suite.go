package testutils

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/nvsco/wallet/pkg/scope"
	prometheusmetrics "github.com/nvsco/wallet/tools/prometheus-metrics"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/suite"
)

type ServiceTestSuite struct {
	suite.Suite

	ServiceName string
	Logger      *logrus.Entry
}

func (s *ServiceTestSuite) SetupSuite() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	logrus.SetReportCaller(true)

	if _, exists := os.LookupEnv("CI"); exists {
		// disable logs with the level below "error" (info/warning/debug)
		logrus.SetLevel(logrus.ErrorLevel)
	}

	logger := logrus.WithField("service", s.ServiceName)

	logger.Info("starting")

	// nolint:errcheck
	_ = ExportEnvForTestsOnce()

	s.Logger = logger
}

func (s *ServiceTestSuite) NewScopeWithTimeout(t time.Duration) (*scope.Scope, context.CancelFunc) {
	ctx, cancel := context.WithTimeout(context.Background(), t)

	xCorrelationID := fmt.Sprintf("test-%s", uuid.New().String())
	return scope.NewScope(ctx, s.Logger, xCorrelationID, prometheusmetrics.NewMetrics()), cancel
}

func (s *ServiceTestSuite) NewScope() *scope.Scope {
	xCorrelationID := fmt.Sprintf("test-%s", uuid.New().String())
	return scope.NewScope(context.Background(), s.Logger, xCorrelationID, prometheusmetrics.NewMetrics())
}

func (s *ServiceTestSuite) RequireJSONEquals(in, out interface{}) {
	inJSON, err := json.Marshal(in)
	s.Require().NoError(err)

	outJSON, err := json.Marshal(out)
	s.Require().NoError(err)

	s.Require().JSONEq(string(inJSON), string(outJSON))
}
