package scrape

import (
	"github.com/tvanriel/price-tracker/internal/curlhttp"
	"github.com/tvanriel/price-tracker/internal/extractors"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("scrape",
	fx.Provide(
		FxNewScraper,
	),
)

type FxNewScraperParams struct {
	fx.In
	Log       *zap.Logger
	Config    Config
	Extractor extractors.Extractor
	CurlHTTP  *curlhttp.CurlHttp
}

func FxNewScraper(p FxNewScraperParams) *Scraper {
	return NewScraper(NewScraperParams{
		Log:        p.Log,
		Config:     p.Config,
		Extractors: p.Extractor,
		CurlHTTP:   p.CurlHTTP,
	})
}
