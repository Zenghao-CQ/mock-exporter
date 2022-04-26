/**
 *  Author: SongLee24
 *  Email: lisong.shine@qq.com
 *  Date: 2018-08-15
 *
 *
 *  prometheus.Desc是指标的描述符，用于实现对指标的管理
 *
 */

package collector

import (
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

// 指标结构体
type Metrics struct {
	metrics      map[string]*prometheus.Desc
	mutex        sync.Mutex
	mockLabels   map[string]string
	failureTypes int
	podNames     []string
}

/**
 * 函数：newGlobalMetric
 * 功能：创建指标描述符
 */
func newGlobalMetric(namespace string, metricName string, docString string, labels []string) *prometheus.Desc {
	return prometheus.NewDesc(namespace+"_"+metricName, docString, labels, nil)
}

/**
 * 工厂方法：NewMetrics
 * 功能：初始化指标信息，即Metrics结构体
 */
func NewMetrics(namespace string, failureTypes int) *Metrics {
	return &Metrics{
		metrics: map[string]*prometheus.Desc{
			"failure_metric": newGlobalMetric(namespace, "failure_counter_metric", "The description of failure an fuilure", []string{"type", "time", "timeFormat", "visited", "podName"}),
			// "my_gauge_metric":   newGlobalMetric(namespace, "my_gauge_metric", "The description of my_gauge_metric", []string{"host"}),
		},
		mockLabels: map[string]string{
			"type":       "0",
			"time":       strconv.FormatInt(time.Now().UnixNano()/1e6, 10),
			"timeFormat": time.Now().In(time.FixedZone("UTC+8", 8*60*60)).Format("2006-01-02 15:04:05"),
			"visited":    "false",
			"podName":    "busybox",
		},
		failureTypes: failureTypes,
	}
}

/**
 * 接口：Describe
 * 功能：传递结构体中的指标描述符到channel
 */
func (c *Metrics) Describe(ch chan<- *prometheus.Desc) {
	for _, m := range c.metrics {
		ch <- m
	}
}

/**
 * 接口：Collect
 * 功能：抓取最新的数据，传递给channel
 */
func (c *Metrics) Collect(ch chan<- prometheus.Metric) {
	c.mutex.Lock() // 加锁
	defer c.mutex.Unlock()
	v, _ := strconv.Atoi(c.mockLabels["type"])
	ch <- prometheus.MustNewConstMetric(c.metrics["failure_metric"], prometheus.CounterValue, float64(v), c.mockLabels["type"], c.mockLabels["time"], c.mockLabels["timeFormat"], c.mockLabels["visited"], c.mockLabels["podName"])
	c.mockLabels["visited"] = "true"
}

/**
 * 函数：GenerateMockData
 * 功能：生成模拟数据
 */
func (c *Metrics) GenerateMockData() {
	c.mutex.Lock() // 加锁
	defer c.mutex.Unlock()
	ct := time.Now()
	c.mockLabels["type"] = strconv.Itoa(rand.Intn(c.failureTypes) + 1)
	c.mockLabels["time"] = strconv.FormatInt(int64(ct.Nanosecond())/1e3, 10)
	c.mockLabels["timeFormat"] = ct.In(time.FixedZone("UTC+8", 8*60*60)).Format("2006-01-02 15:04:05.000")
	c.mockLabels["visited"] = "false"
	c.mockLabels["podName"] = c.podNames[rand.Int()%len(c.podNames)]
	log.Printf("Generate failure type:%s time:%s timesamp:%s visited:%s podName:%s", c.mockLabels["type"], c.mockLabels["time"], c.mockLabels["timeFormat"], c.mockLabels["visited"], c.mockLabels["podName"])
}
