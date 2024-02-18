package metrics

import (
	"sync"

	pprom "github.com/prometheus/client_golang/prometheus"
	"github.com/tvanriel/cloudsdk/prometheus"
	"github.com/tvanriel/price-tracker/internal/scrape"
	"go.uber.org/fx"
)

type NewMetricsParams struct {
	fx.In

	Prometheus *prometheus.Prometheus
}
type Metrics struct {
	Prometheus       *prometheus.Prometheus
	PriceDescs       map[string]*pprom.Desc
	Prices           map[string]float64
	RequestTimeDescs map[string]*pprom.Desc
	RequestTimes     map[string]float64
	BufferTimeDescs  map[string]*pprom.Desc
	BufferTimes      map[string]float64
	ScriptTimeDescs  map[string]*pprom.Desc
	ScriptTimes      map[string]float64
	TotalTimeDescs   map[string]*pprom.Desc
	TotalTimes       map[string]float64
	m                sync.RWMutex
}

func NewMetrics(params NewMetricsParams, lc fx.Lifecycle) *Metrics {
	return &Metrics{
		Prometheus:       params.Prometheus,
		PriceDescs:       make(map[string]*pprom.Desc),
		Prices:           make(map[string]float64),
		RequestTimeDescs: make(map[string]*pprom.Desc),
		RequestTimes:     make(map[string]float64),
		BufferTimeDescs:  make(map[string]*pprom.Desc),
		BufferTimes:      make(map[string]float64),
		ScriptTimeDescs:  make(map[string]*pprom.Desc),
		ScriptTimes:      make(map[string]float64),
		TotalTimeDescs:   make(map[string]*pprom.Desc),
		TotalTimes:       make(map[string]float64),
		m:                sync.RWMutex{},
	}
}

type RegisterScrapeMetricsParams struct {
	Name   string
	Labels map[string]string
}

func (m *Metrics) AddDescription(r *RegisterScrapeMetricsParams) {

	labels := make(map[string]string)
	for k, v := range r.Labels {
		labels[k] = v
	}
	labels["name"] = r.Name

	m.PriceDescs[r.Name] = pprom.NewDesc("price", "Price of the article", nil, labels)
}
func (m *Metrics) Describe(ch chan<- *pprom.Desc) {
	for k := range m.PriceDescs {
		ch <- m.PriceDescs[k]
	}
	for k := range m.RequestTimeDescs {
		ch <- m.RequestTimeDescs[k]
	}
	for k := range m.BufferTimeDescs {
		ch <- m.BufferTimeDescs[k]
	}
	for k := range m.ScriptTimeDescs {
		ch <- m.ScriptTimeDescs[k]
	}
	for k := range m.TotalTimeDescs {
		ch <- m.TotalTimeDescs[k]
	}

}

func (m *Metrics) Collect(ch chan<- pprom.Metric) {
	m.m.RLock()
	for k := range m.PriceDescs {
		ch <- pprom.MustNewConstMetric(
			m.PriceDescs[k],
			pprom.GaugeValue,
			m.Prices[k],
		)
	}
	for k := range m.RequestTimeDescs {
		ch <- pprom.MustNewConstMetric(
			m.RequestTimeDescs[k],
			pprom.GaugeValue,
			m.RequestTimes[k],
		)
	}
	for k := range m.BufferTimeDescs {
		ch <- pprom.MustNewConstMetric(
			m.BufferTimeDescs[k],
			pprom.GaugeValue,
			m.BufferTimes[k],
		)
	}
	for k := range m.ScriptTimeDescs {
		ch <- pprom.MustNewConstMetric(
			m.ScriptTimeDescs[k],
			pprom.GaugeValue,
			m.ScriptTimes[k],
		)
	}
	for k := range m.TotalTimeDescs {
		ch <- pprom.MustNewConstMetric(
			m.TotalTimeDescs[k],
			pprom.GaugeValue,
			m.TotalTimes[k],
		)
	}
	m.m.RUnlock()

}

func (m *Metrics) Record(res *scrape.ScrapeResult) {
	m.m.Lock()
	m.Prices[res.Name] = (res.Price)
	//m.RequestTimes[res.Name] = (float64(res.RequestTime))
	//m.BufferTimes[res.Name] = (float64(res.BufferTime))
	//m.ScriptTimes[res.Name] = (float64(res.ScriptTime))
	//m.TotalTimes[res.Name] = (float64(res.TotalTime))
	m.m.Unlock()
}
