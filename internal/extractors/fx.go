package extractors

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("extractors", fx.Provide(FxNewTengoExtractor))

type FxNewTengoExtractorParams struct {
	fx.In

	Config Config
	Log    *zap.Logger
}

func FxNewTengoExtractor(p FxNewTengoExtractorParams) Extractor {
	return NewTengoExtractors(NewTengoExtractorParams{
		Config: p.Config,
		Log:    p.Log,
	})
}
