package app

import (
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/cloudsdk/prometheus"
	"github.com/tvanriel/price-tracker/internal/application"
	"github.com/tvanriel/price-tracker/internal/config"
	"github.com/tvanriel/price-tracker/internal/curlhttp"
	"github.com/tvanriel/price-tracker/internal/extractors"
	"github.com/tvanriel/price-tracker/internal/metrics"
	"github.com/tvanriel/price-tracker/internal/scrape"
	"go.uber.org/fx"
)

func Application() {

	fx.New(
		prometheus.Module,
		logging.Module,
		curlhttp.Module,
		extractors.Module,
		logging.FXLogger(),
		config.Configuration(),
		scrape.Module,
		metrics.Module,
		application.Module,
		fx.Invoke(func(_ *application.Application) {}),
	).Run()

}
