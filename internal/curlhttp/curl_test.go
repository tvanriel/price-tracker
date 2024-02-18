package curlhttp_test

import (
	"context"
	"strings"
	"testing"

	"github.com/tvanriel/price-tracker/internal/curlhttp"
	"go.uber.org/zap"
	"gotest.tools/v3/assert"
)

func TestMain(t *testing.T) {
	c := curlhttp.NewCurlHttp(curlhttp.NewCurlHttpParams{Binary: "/usr/bin/curl", Log: zap.L()})
	req := &curlhttp.CurlRequest{}
	req.HTTP2()
	req.EnableCompression()
	req.SetMethod("GET")
	req.SetURL("https://www.aldi.nl/producten/brood-bakkerij/dagvers-brood/wit-rond-3931-1-0.article.html")
	req.SetUserAgent("Price-Tracker v1.0.0 (+https://github.com/tvanriel)")
	req.SetCookies(make(map[string]string))
	req.SetBody(strings.NewReader(""))
	req.AddHeader("Accept", []string{"text/html"})

	resp, err := c.Do(context.Background(), req)

	assert.NilError(t, err)
	assert.Assert(t, len(resp.Headers()) == 20)
	assert.Equal(t, resp.HttpCode(), "200")
}
