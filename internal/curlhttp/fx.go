package curlhttp

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module("curl", fx.Provide(FxNewCurlHttp))

type FxCurlHttpParams struct {
	fx.In

	Config Config
	Log    *zap.Logger
}

func FxNewCurlHttp(p FxCurlHttpParams) *CurlHttp {
	return NewCurlHttp(NewCurlHttpParams{
		Config: p.Config,
		Log:    p.Log,
	})
}
