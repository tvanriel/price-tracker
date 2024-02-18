package application

import (
	"time"

	"github.com/tvanriel/cloudsdk/prometheus"
	"github.com/tvanriel/price-tracker/internal/metrics"
	"github.com/tvanriel/price-tracker/internal/scrape"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type NewApplicationParams struct {
	fx.In

	Metrics *metrics.Metrics
	Scraper *scrape.Scraper
	Log     *zap.Logger
        Prometheus *prometheus.Prometheus
}
type Application struct {
	Metrics *metrics.Metrics
	Scraper *scrape.Scraper
	Log     *zap.Logger
        Prometheus *prometheus.Prometheus
}

func NewApplication(p NewApplicationParams, lc fx.Lifecycle) *Application {
	app := &Application{
		Scraper: p.Scraper,
		Metrics: p.Metrics,
		Log:     p.Log.Named("app"),
                Prometheus: p.Prometheus,
	}
	lc.Append(fx.StartHook(func() error {
		go func() {
			ch := make(chan *scrape.ScrapeResult)
			go app.DescribeMetrics()
			go app.StartScraping(ch)
			go app.PipeMetrics(ch)
		}()
		return nil
	}))
	return app
}

func (a *Application) DescribeMetrics() {
	targets := a.Scraper.Config.Targets
	for t := range targets {

		a.Metrics.AddDescription(&metrics.RegisterScrapeMetricsParams{
			Name:   targets[t].Name,
			Labels: targets[t].Labels,
		})
	}
        a.Prometheus.Register(a.Metrics)
        
}
func (a *Application) StartScraping(ch chan *scrape.ScrapeResult) {
	a.Scraper.Scrape(ch)
	ticker := time.NewTicker(3 * time.Hour)
	for _ = range ticker.C {
		a.Scraper.Scrape(ch)
	}

}
func (a *Application) PipeMetrics(ch chan *scrape.ScrapeResult) {
	for res := range ch {
		a.Metrics.Record(res)
	}
}
