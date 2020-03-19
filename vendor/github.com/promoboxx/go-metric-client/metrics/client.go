package metrics

import (
	"context"
	"time"

	"github.com/opentracing/opentracing-go"
)

// Client is the interface that is used for communicating with a metrics
// aggregator. It defines the functions for each of the different types of
// metrics we can interact with

//go:generate mockgen -destination=./mockmetrics/mock-metrics.go -package=mockmetrics github.com/promoboxx/go-metric-client/metrics Client

type Client interface {

	// This metric is used to count how often background jobs run. This should be automatically created and incremented for all go background jobs (sessions)
	BackgroundRate(sessionID, jobName string, params map[string]string, value int64) error

	// This metric is used to count how often background jobs error. This should be incremented automatically for all go background jobs (sessions) that end due to an error.
	BackgroundError(sessionID, jobName string, params map[string]string, code, message string, value int64) error

	// This gauge metric is used to keep track of the runtime of various jobs. This should be automatically tracked and submitted for all go background jobs (sessions), regardless of if they end in an error or in success.
	BackgroundDuration(sessionID, jobName string, params map[string]string, value time.Duration) error

	// This metric is used to keep track of business process counters in background jobs (sessions). They should be manually invoked whenever there is a custom thing we need to track
	BackgroundCustom(sessionID string, jobName string, customName string, params, other map[string]string, value int64) error

	// This metric is used to count how often we communicate with an external partner we are integrated with. This communication could either be initiated by us within a service or could be to a webhook
	// we have set up for that external service to reach out to us. This should be manually added to all partner-integrated code paths, however, if we have partner-specific integration packages,
	//these metrics should be instrumented at that level.
	ExternalRate(direction, externalService, path string, value int64) error

	// This metric is used to count how often partner communications error. This should be incremented automatically for all partner communications that end due to an error.
	ExternalError(direction, externalService, path, code, message string, value int64) error

	// This gauge metric is used to keep track of the runtime of various partner communications. This should be automatically tracked and submitted for all partner integration communications,
	// regardless of if they end in an error or in success.
	ExternalDuration(direction, externalService, path string, value time.Duration) error

	// This metric is used to keep track of business process counters in partner communications. They should be manually invoked whenever there is a custom thing we need to track
	ExternalCustom(direction, externalService, path, customName string, other map[string]string, value int64) error

	// This metric is used to keep track of business process counters in internal communications. They should be manually invoked whenever there is a custom thing we need to track. For instance,
	// the ad service might keep track of the number of ad failures related to the wallet service.
	InternalCustom(originatingService, destinationService, path, customName string, other map[string]string, value int64) error

	// This is to allow the metric client to fulfil the Tracer interface within the go-session-lock package
	StartSpanWithContext(ctx context.Context, name string) (opentracing.Span, context.Context)
}
