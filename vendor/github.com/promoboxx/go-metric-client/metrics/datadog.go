package metrics

import (
	"context"
	"sync"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/opentracing/opentracing-go"
)

var instance *datadogMetricsClient
var once sync.Once
var metricErr error

// struct for datadogMetricsClient
type datadogMetricsClient struct {
	statsdClient statsd.ClientInterface
}

// NewDatadogMetricsClient returns a new metrics client that implements the Client interface
// defined by this package which also works as a singleton to control object creation
func NewDatadogMetricsClient(address string, options statsd.Option, service string, baseTag map[string]string) (Client, error) {

	once.Do(func() {
		newBaseTag := []string{"service:" + service}

		for key, val := range baseTag {
			newBaseTag = append(newBaseTag, key+":"+val)
		}

		statsdClient, err := statsd.New(address, statsd.WithTags(newBaseTag))
		if err != nil {
			metricErr = err
		}

		instance = &datadogMetricsClient{statsdClient: statsdClient}
	})

	return instance, metricErr
}

// This metric is used to count how often background jobs run.
func (dmc *datadogMetricsClient) BackgroundRate(sessionID, jobName string, params map[string]string, value int64) error {

	metricTag := []string{"session_id:" + sessionID, "job_name:" + jobName}

	sanitizedMetricTag := tagsBuilder(metricTag, params, nil)

	return dmc.statsdClient.Count("pbxx.background.rate", value, sanitizedMetricTag, 0)
}

// This metric is used to count how often background jobs error.
func (dmc *datadogMetricsClient) BackgroundError(sessionID, jobName string, params map[string]string, code, message string, value int64) error {

	metricTag := []string{"session_id:" + sessionID, "job_name:" + jobName, "code:" + code, "message:" + message}

	sanitizedMetricTag := tagsBuilder(metricTag, params, nil)

	return dmc.statsdClient.Count("pbxx.background.error", value, sanitizedMetricTag, 0)

}

// This gauge metric is used to keep track of the runtime of various jobs.
func (dmc *datadogMetricsClient) BackgroundDuration(sessionID, jobName string, params map[string]string, value time.Duration) error {

	metricTag := []string{"session_id:" + sessionID, "job_name:" + jobName}

	sanitizedMetricTag := tagsBuilder(metricTag, params, nil)

	return dmc.statsdClient.Timing("pbxx.background.duration", value, sanitizedMetricTag, 0)
}

// This metric is used to keep track of business process counters in background jobs (sessions).
func (dmc *datadogMetricsClient) BackgroundCustom(sessionID, jobName, customName string, params, other map[string]string, value int64) error {

	metricTag := []string{"session_id:" + sessionID, "job_name:" + jobName, "custom_name:" + customName}

	sanitizedMetricTag := tagsBuilder(metricTag, params, other)

	return dmc.statsdClient.Count("pbxx.background.custom", value, sanitizedMetricTag, 0)
}

// This metric is used to count how often we communicate with an external partner we are integrated with.
func (dmc *datadogMetricsClient) ExternalRate(direction, externalService, path string, value int64) error {
	metricTag := sanitizeTags([]string{"direction:" + direction, "external_service:" + externalService, "path:" + path})

	return dmc.statsdClient.Count("pbxx.external.rate", value, metricTag, 0)
}

// This metric is used to count how often partner communications error.
func (dmc *datadogMetricsClient) ExternalError(direction, externalService, path, code, message string, value int64) error {
	metricTag := sanitizeTags([]string{"direction:" + direction, "external_service:" + externalService, "path:" + path, "code:" + code, "message:" + message})

	return dmc.statsdClient.Count("pbxx.external.error", value, metricTag, 0)
}

// This gauge metric is used to keep track of the runtime of various partner communications.
func (dmc *datadogMetricsClient) ExternalDuration(direction, externalService, path string, value time.Duration) error {
	metricTag := sanitizeTags([]string{"direction:" + direction, "external_service:" + externalService, "path:" + path})

	return dmc.statsdClient.Timing("pbxx.external.duration", value, metricTag, 0)
}

// This metric is used to keep track of business process counters in partner communications.
func (dmc *datadogMetricsClient) ExternalCustom(direction, externalService, path, customName string, other map[string]string, value int64) error {
	metricTag := []string{"direction:" + direction, "external_service:" + externalService, "path:" + path, "custom_name:" + customName}

	sanitizedMetricTag := tagsBuilder(metricTag, nil, other)

	return dmc.statsdClient.Count("pbxx.external.custom", value, sanitizedMetricTag, 0)
}

// This metric is used to keep track of business process counters in internal communications.
func (dmc *datadogMetricsClient) InternalCustom(originatingService, destinationService, path, customName string, other map[string]string, value int64) error {
	metricTag := []string{"originating_service:" + originatingService, "destination_service:" + destinationService, "path:" + path, "custom_name:" + customName}

	sanitizedMetricTag := tagsBuilder(metricTag, nil, other)

	return dmc.statsdClient.Count("pbxx.internal.custom", value, sanitizedMetricTag, 0)
}

// This metric is used to keep track of business process counters in internal communications.
func (dmc *datadogMetricsClient) StartSpanWithContext(ctx context.Context, name string) (opentracing.Span, context.Context) {

	span, ctx := opentracing.StartSpanFromContext(ctx, name)
	return span, opentracing.ContextWithSpan(ctx, span)
}
