package collector

import (
	"fmt"
	"io/ioutil"
	"context"
	"log"
	"os"
	"net/http"
	"sync"
	"time"
	// "reflect"

	"github.com/didi/nightingale/src/dataobj"
	"github.com/mainto-hz/n9e-query-prometheus/config"
	"github.com/prometheus/client_golang/api"
	"github.com/prometheus/client_golang/api/prometheus/v1"

)

func Gather() []*dataobj.MetricValue {
	var wg sync.WaitGroup
	var res []*dataobj.MetricValue

	cfg := config.Get()
	metricChan := make(chan *dataobj.MetricValue)
	done := make(chan struct{}, 1)

	go func() {
		// fmt.Println("123123123123")
		// fmt.Printf("asdfasdfasdfasdf $v\n", 1)
		defer func() { done <- struct{}{} }()
		for m := range metricChan {
			res = append(res, m)
		}
		// fmt.Printf("nid: %v\n", cfg.Nid)

		// if cfg.Nid != ""{
		// 	for m := range metricChan {
		// 		res.Nid = cfg.Nid
		// 	}
		// }
	}()

	for _, exporterUrl := range cfg.ExporterUrls {
		wg.Add(1)
		go func(url string,query string) {
			defer wg.Done()
			if metrics, err := gatherExporter(url,query); err == nil {
				for _, m := range metrics {
					if typ, exists := cfg.MetricType[m.Metric]; exists {
						m.CounterType = typ
					}

					if cfg.MetricPrefix != "" {
						m.Metric = cfg.MetricPrefix + m.Metric
					}
					metricChan <- m
				}
			}
		}(exporterUrl,cfg.Query)
	}

	wg.Wait()
	close(metricChan)

	<-done

	return res
}

func gatherExporter(url string, query string) ([]*dataobj.MetricValue, error) {
	// body, err := gatherExporterUrl(url)
	// if err != nil {
	// 	log.Printf("gather metrics from exporter error, url :[%s] ,error :%v", url, err)
	// 	return nil, err
	// }

	// metrics, err := Parse(body)
	// if err != nil {
	// 	log.Printf("parse metrics error, url :[%s] ,error :%v", url, err)
	// 	return nil, err
	// }

	// return metrics, nil


	client, err := api.NewClient(api.Config{
		Address: url,
	})
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		os.Exit(1)
	}

	v1api := v1.NewAPI(client)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	result, warnings, err := v1api.Query(ctx, query, time.Now())
	if err != nil {
		fmt.Printf("Error querying Prometheus: %v\n", err)
		os.Exit(1)
	}
	if len(warnings) > 0 {
		fmt.Printf("Warnings: %v\n", warnings)
	}
	// fmt.Printf("Result:\n%v\n", result)
	// fmt.Println(reflect.TypeOf(result))
	// rs := result.String()
	// fmt.Printf("Result:\n%v\n", rs)
	// fmt.Println(reflect.TypeOf(rs))

	// metrics, err := Parse([]byte(rs))
	metrics, err := Parse(result)
	if err != nil {
		log.Printf("parse metrics error ,error :%v", err)
		return nil, err
	}

	return metrics, nil

	// return nil,nil
}

func gatherExporterUrl(url string) ([]byte, error) {
	var buf []byte
	var req *http.Request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return buf, err
	}

	client := &http.Client{
		Timeout: time.Duration(config.Get().Timeout) * time.Millisecond,
	}

	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return buf, fmt.Errorf("error making HTTP request to %s: %s", url, err)
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return buf, fmt.Errorf("%s returned HTTP status %s", url, resp.Status)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return buf, fmt.Errorf("error reading body: %s", err)
	}

	return body, nil
}
