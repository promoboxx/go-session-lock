package metrics

import (
	"context"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/opentracing/opentracing-go"
	"github.com/sirupsen/logrus"
)

// struct for the asyncDatadogClient
type asyncDatadogClient struct {
	dmc    Client
	logger *logrus.Entry
}

// NewAsyncDatadogMetricsClient returns a new metrics client that implements the Client interface
// defined by this package which also works as a singleton to control object creation. This allows the metrics
//  to be run concurrently through Goroutines through wrapping.
func NewAsyncDatadogClient(address string, options statsd.Option, service string, baseTag map[string]string, logger *logrus.Entry) (Client, error) {
	datadogClient, err := NewDatadogMetricsClient(address, options, service, baseTag)

	if err != nil {
		return nil, err
	}

	return &asyncDatadogClient{datadogClient, logger}, nil
}

// This metric is used to count how often background jobs run.
func (a asyncDatadogClient) BackgroundRate(sessionID, jobName string, params map[string]string, value int64) error {
	go func() {
		err := a.BackgroundRate(sessionID, jobName, params, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

// This metric is used to count how often background jobs error.
func (a asyncDatadogClient) BackgroundError(sessionID, jobName string, params map[string]string, code, message string, value int64) error {
	go func() {
		err := a.BackgroundError(sessionID, jobName, params, code, message, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

// This gauge metric is used to keep track of the runtime of various jobs.
func (a asyncDatadogClient) BackgroundDuration(sessionID, jobName string, params map[string]string, value time.Duration) error {
	go func() {
		err := a.BackgroundDuration(sessionID, jobName, params, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

// This metric is used to keep track of business process counters in background jobs (sessions).
func (a asyncDatadogClient) BackgroundCustom(sessionID string, jobName string, customName string, params, other map[string]string, value int64) error {
	go func() {
		err := a.BackgroundCustom(sessionID, jobName, customName, params, other, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

// This metric is used to count how often we communicate with an external partner we are integrated with.
func (a asyncDatadogClient) ExternalRate(direction, externalService, path string, value int64) error {
	go func() {
		err := a.ExternalRate(direction, externalService, path, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

// This metric is used to count how often partner communications error.
func (a asyncDatadogClient) ExternalError(direction, externalService, path, code, message string, value int64) error {
	go func() {
		err := a.ExternalError(direction, externalService, path, code, message, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

// This gauge metric is used to keep track of the runtime of various partner communications.
func (a asyncDatadogClient) ExternalDuration(direction, externalService, path string, value time.Duration) error {
	go func() {
		err := a.ExternalDuration(direction, externalService, path, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

// This metric is used to keep track of business process counters in partner communications.
func (a asyncDatadogClient) ExternalCustom(direction, externalService, path, customName string, other map[string]string, value int64) error {
	go func() {
		err := a.ExternalCustom(direction, externalService, path, customName, other, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

// This metric is used to keep track of business process counters in internal communications.
func (a asyncDatadogClient) InternalCustom(originatingService, destinationService, path, customName string, other map[string]string, value int64) error {
	go func() {
		err := a.InternalCustom(originatingService, destinationService, path, customName, other, value)

		if err != nil {
			a.logger.Errorf("error sending metrics data: %s", err)
		}
	}()
	return nil
}

func (a asyncDatadogClient) StartSpanWithContext(ctx context.Context, name string) (opentracing.Span, context.Context) {

	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	return span, opentracing.ContextWithSpan(ctx, span)
}
