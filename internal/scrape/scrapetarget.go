package scrape

import (
	rand "math/rand/v2"
	"time"

	"github.com/tvanriel/price-tracker/internal/curlhttp"
	"github.com/tvanriel/price-tracker/internal/extractors"
	"go.uber.org/zap"
	"golang.org/x/net/context"
)

type NewScraperParams struct {
	Log        *zap.Logger
	Config     Config
	Extractors extractors.Extractor
	CurlHTTP   *curlhttp.CurlHttp
}

type Scraper struct {
	Log        *zap.Logger
	Config     Config
	Extractors extractors.Extractor
	CurlHTTP   *curlhttp.CurlHttp
}

func NewScraper(p NewScraperParams) *Scraper {
	return &Scraper{
		Log:        p.Log.Named("scraper"),
		Config:     p.Config,
		Extractors: p.Extractors,
		CurlHTTP:   p.CurlHTTP,
	}
}

func (s *Scraper) Scrape(ch chan *ScrapeResult) {

	for t := range s.Config.Targets {
		go func() {
                        userAgent := s.Config.UserAgents[rand.IntN(len(s.Config.UserAgents))]
			name := s.Config.Targets[t].Name

			req := curlRequestFromTarget(&s.Config.Targets[t], userAgent)
                        s.Log.Info("Requesting", zap.String("name", name))
                        ctx, _ := context.WithTimeout(context.Background(), time.Second * 45)
			resp, err := s.CurlHTTP.Do(ctx, req)
			if err != nil {
				s.Log.Error("Cannot perform request", zap.String("name", name), zap.Error(err))
				return
			}
			price, err := s.Extractors.Get(resp.Body(), s.Config.Targets[t].Script)
			if err != nil {
				s.Log.Error("Cannot scrape target", zap.String("name", name), zap.Error(err))
				return
			}
			res := &ScrapeResult{
				Name:   name,
				Labels: s.Config.Targets[t].Labels,
				Price:  price,
			}
			ch <- res
		}()
	}
}

type ScrapeResult struct {
	Name   string
	Labels map[string]string
	Price  float64
}

func curlRequestFromTarget(t *Target, userAgent string) *curlhttp.CurlRequest {
	req := &curlhttp.CurlRequest{}
	req.HTTP1()
	req.SetMethod("GET")
	req.SetUserAgent(userAgent)
	if t.HTTPVersion == "2" {
		req.HTTP2()
	}
	if t.Compression {
		req.EnableCompression()
	}

	req.SetURL(t.Page)
	req.SetCookies(t.Cookies)
	for header := range t.Headers {
		req.AddHeader(header, t.Headers[header])
	}
	return req

}
