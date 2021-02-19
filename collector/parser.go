package collector

import (
	// "bufio"
	// "bytes"
	// "fmt"
	// "math"
	// "strings"
	"time"
	// "reflect"

	"github.com/didi/nightingale/src/dataobj"
	// dto "github.com/prometheus/client_model/go"
	// "github.com/prometheus/common/expfmt"
	promodel "github.com/prometheus/common/model"

	"github.com/mainto-hz/n9e-query-prometheus/config"
	"github.com/mainto-hz/n9e-query-prometheus/model"
)

var now = time.Now().Unix()

// Get labels from metric
func makeLabels(m promodel.Metric) map[string]string {
	tags := map[string]string{}
	for name, val := range m {
		if name == "endpoint"{
			continue
		}
		tags[string(name)] = string(val)
	}
	return tags
}
// func Parse(buf []byte) ([]*dataobj.MetricValue, error) {
func Parse(buf promodel.Value) ([]*dataobj.MetricValue, error) {
	cfg := config.Get()
	var metricList []*dataobj.MetricValue
	// 数据格式转换，从prometheus的格式转换为json
	for _, data := range buf.(promodel.Vector) {
		
		// fmt.Printf("Metric = %v \n", data.Metric)
		// fmt.Printf("MetricName = %v \n", data.Metric["__name__"])
		// fmt.Printf("Value = %v \n", data.Value)
		// fmt.Printf("Timestamp = %v \n", data.Timestamp)

		now = time.Now().Unix()
		tags := makeLabels(data.Metric)
		
		metric := model.NewGaugeMetric(string(data.Metric["__name__"]), float64(data.Value), now, tags)
		

		
		metric.Endpoint = cfg.Endpoint
		// metric.Nid = ""
		metric.Tags = makeAppendTags(metric.TagsMap, config.AppendTagsMap)
		// set provided Time, ms to s
		// if metric.Timestamp  > 0 {
		// 	metric.Timestamp = metric.Timestamp  / 1000
		// }
		metricList = append(metricList, metric)
	}
	
	
	return metricList, nil
	// -----------------------------------------------------------------------------------------------------------------//

	// var metricList []*dataobj.MetricValue
	// var parser expfmt.TextParser
	// cfg := config.Get()
	// // parse even if the buffer begins with a newline
	// buf = bytes.TrimPrefix(buf, []byte("\n"))
	// // Read raw data
	// buffer := bytes.NewBuffer(buf)
	// reader := bufio.NewReader(buffer)

	// // Prepare output
	// metricFamilies := make(map[string]*dto.MetricFamily)
	// metricFamilies, err := parser.TextToMetricFamilies(reader)
	// if err != nil {
	// 	return nil, fmt.Errorf("reading text format failed: %s", err)
	// }
	// now = time.Now().Unix()
	// // read metrics
	// for basename, mf := range metricFamilies {
	// 	metrics := []*dataobj.MetricValue{}
	// 	for _, m := range mf.Metric {
	// 		// pass ignore metric
	// 		if filterIgnoreMetric(basename) {
	// 			continue
	// 		}
	// 		switch mf.GetType() {
	// 		case dto.MetricType_GAUGE:
	// 			// gauge metric
	// 			metrics = makeCommon(basename, m)
	// 		case dto.MetricType_COUNTER:
	// 			// counter metric
	// 			metrics = makeCommon(basename, m)
	// 		case dto.MetricType_SUMMARY:
	// 			// summary metric
	// 			metrics = makeQuantiles(basename, m)
	// 		case dto.MetricType_HISTOGRAM:
	// 			// histogram metric
	// 			metrics = makeBuckets(basename, m)
	// 		case dto.MetricType_UNTYPED:
	// 			// untyped as gauge
	// 			metrics = makeCommon(basename, m)
	// 		}

	// 		// render endpoint info
	// 		for _, metric := range metrics {
	// 			metric.Endpoint = cfg.Endpoint
	// 			metric.Tags = makeAppendTags(metric.TagsMap, config.AppendTagsMap)
	// 			// set provided Time, ms to s
	// 			if m.GetTimestampMs() > 0 {
	// 				metric.Timestamp = m.GetTimestampMs() / 1000
	// 			}
	// 			metricList = append(metricList, metric)
	// 		}
	// 	}
	// }

	// return metricList, err
}

// Get Quantiles from summary metric
// func makeQuantiles(basename string, m *dto.Metric) []*dataobj.MetricValue {
// 	metrics := []*dataobj.MetricValue{}
// 	tags := makeLabels(m)

// 	countName := fmt.Sprintf("%s_count", basename)
// 	metrics = append(metrics, model.NewCumulativeMetric(countName, m.GetSummary().SampleCount, now, tags))

// 	sumName := fmt.Sprintf("%s_sum", basename)
// 	metrics = append(metrics, model.NewCumulativeMetric(sumName, m.GetSummary().SampleSum, now, tags))

// 	for _, q := range m.GetSummary().Quantile {
// 		tagsNew := make(map[string]string)
// 		for tagKey, tagValue := range tags {
// 			tagsNew[tagKey] = tagValue
// 		}
// 		if !math.IsNaN(q.GetValue()) {
// 			tagsNew["quantile"] = fmt.Sprint(q.GetQuantile())

// 			metrics = append(metrics, model.NewGaugeMetric(basename, float64(q.GetValue()), now, tagsNew))
// 		}
// 	}

// 	return metrics
// }

// Get Buckets from histogram metric
// func makeBuckets(basename string, m *dto.Metric) []*dataobj.MetricValue {
// 	metrics := []*dataobj.MetricValue{}
// 	tags := makeLabels(m)

// 	countName := fmt.Sprintf("%s_count", basename)
// 	metrics = append(metrics, model.NewCumulativeMetric(countName, m.GetHistogram().SampleCount, now, tags))

// 	sumName := fmt.Sprintf("%s_sum", basename)
// 	metrics = append(metrics, model.NewCumulativeMetric(sumName, m.GetHistogram().SampleSum, now, tags))

// 	for _, b := range m.GetHistogram().Bucket {
// 		tagsNew := make(map[string]string)
// 		for tagKey, tagValue := range tags {
// 			tagsNew[tagKey] = tagValue
// 		}
// 		tagsNew["le"] = fmt.Sprint(b.GetUpperBound())

// 		bucketName := fmt.Sprintf("%s_bucket", basename)
// 		metrics = append(metrics, model.NewGaugeMetric(bucketName, float64(b.GetCumulativeCount()), now, tagsNew))
// 	}

// 	return metrics
// }

// Get gauge and counter from metric
// func makeCommon(metricName string, m *promodel.Metric) []*dataobj.MetricValue {
// 	var val float64
// 	metrics := []*dataobj.MetricValue{}
// 	tags := makeLabels(m)
// 	if m.Gauge != nil {
// 		if !math.IsNaN(m.GetGauge().GetValue()) {
// 			val = float64(m.GetGauge().GetValue())
// 			metrics = append(metrics, model.NewGaugeMetric(metricName, val, now, tags))
// 		}
// 	} else if m.Counter != nil {
// 		if !math.IsNaN(m.GetCounter().GetValue()) {
// 			val = float64(m.GetCounter().GetValue())
// 			metrics = append(metrics, model.NewCumulativeMetric(metricName, val, now, tags))
// 		}
// 	} else if m.Untyped != nil {
// 		// untyped as gauge
// 		if !math.IsNaN(m.GetUntyped().GetValue()) {
// 			val = float64(m.GetUntyped().GetValue())
// 			metrics = append(metrics, model.NewGaugeMetric(metricName, val, now, tags))
// 		}
// 	}
// 	return metrics
// }



// append tags
func makeAppendTags(tagsMap map[string]string, appendTagsMap map[string]string) string {
	if len(tagsMap) == 0 && len(appendTagsMap) == 0 {
		return ""
	}

	if len(tagsMap) == 0 {
		return dataobj.SortedTags(appendTagsMap)
	}

	if len(appendTagsMap) == 0 {
		return dataobj.SortedTags(tagsMap)
	}

	for k, v := range appendTagsMap {
		if k == "endpoint"{
			continue
		}
		tagsMap[k] = v
	}

	return dataobj.SortedTags(tagsMap)
}

// func filterIgnoreMetric(basename string) bool {
// 	ignorePrefix := config.Get().IgnoreMetricsPrefix
// 	if len(config.Get().IgnoreMetricsPrefix) == 0 {
// 		return false
// 	}

// 	for _, pre := range ignorePrefix {
// 		if strings.HasPrefix(basename, pre) {
// 			return true
// 		}
// 	}
// 	return false
// }
