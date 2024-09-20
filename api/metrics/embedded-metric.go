package metrics

type MetricNamespace struct {
	Namespace  string          `json:"Namespace"`
	Dimensions [][]string      `json:"Dimensions"`
	Metrics    []MetricDetails `json:"Metrics"`
}

type MetricDetails struct {
	Name string `json:"Name"`
	Unit string `json:"Unit"`
}

type EmbeddedMetric struct {
	Aws struct {
		Timestamp         int64             `json:"Timestamp"`
		CloudWatchMetrics []MetricNamespace `json:"CloudWatchMetrics"`
	} `json:"_aws"`
	OperationName string `json:"OperationName"`
	StatusCode    int    `json:"StatusCode"`
	Latency       int64  `json:"Latency"`
	CallCount     int    `json:"CallCount"`
}
