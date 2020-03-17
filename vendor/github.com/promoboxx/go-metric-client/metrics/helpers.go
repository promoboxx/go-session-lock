package metrics

import (
	"regexp"
	"strings"
)

var whitelist *regexp.Regexp

func init() {
	whitelist = regexp.MustCompile("[^a-z.:0-9\\\\_/-]")
}

// Helper method to build parameters into the metric tag
func tagsBuilder(metricTag []string, params map[string]string, custom map[string]string) []string {

	for key, val := range params {
		metricTag = append(metricTag, "param_"+strings.ToLower(key)+":"+strings.ToLower(val))
	}

	for key, val := range custom {
		metricTag = append(metricTag, "other_"+strings.ToLower(key)+":"+strings.ToLower(val))
	}
	return sanitizeTags(metricTag)
}

// Helper method to sanitize tags that to make sure that the given tags are allowed by Datadog metrics as init in the init whitelist regex above as well as https://docs.datadoghq.com/tagging/#defining-tags
func sanitizeTags(metricTag []string) []string {
	for i, str := range metricTag {
		metricTag[i] = whitelist.ReplaceAllLiteralString(str, "_")
		if len(str) > 200 {
			metricTag[i] = metricTag[i][0:200]
		}

		metricTag[i] = strings.TrimSuffix(metricTag[i], `:`)
	}
	return metricTag
}
