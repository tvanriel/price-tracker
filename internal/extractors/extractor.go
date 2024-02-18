package extractors

import (
	"io"

	tengo "github.com/d5/tengo/v2"
	"go.uber.org/zap"
)

// Extractor receives a response body in it's entirety and extracts the price from it.
type Extractor interface {
	Get(page io.Reader, script string) (float64, error)
}

type NewTengoExtractorParams struct {
	Config Config
	Log    *zap.Logger
}

type TengoExtractor struct {
	config Config
	Log    *zap.Logger
}

func NewTengoExtractors(p NewTengoExtractorParams) *TengoExtractor {
	return &TengoExtractor{
		config: p.Config,
		Log:    p.Log.Named("tengo_extractor"),
	}
}
func (t *TengoExtractor) AsObject() tengo.Object {
	m := &tengo.Map{Value: make(map[string]tengo.Object)}
	for e := range t.config.InnerText {
		m.IndexSet(
			&tengo.String{Value: e},
			TengoInnerTextExtractor(t.config.InnerText[e].Selector, t.Log.With(zap.String("extractor", e))),
		)
	}
	for e := range t.config.Attribute {
		m.IndexSet(
			&tengo.String{Value: e},
			TengoAttributeExtractor(t.config.Attribute[e].Selector, t.config.Attribute[e].Attribute, t.Log.With(zap.String("extractor", e))),
		)
	}
	return m
}
func (t *TengoExtractor) Get(page io.Reader, script string) (float64, error) {
	s := tengo.NewScript([]byte(script))
	s.Add("ex", t.AsObject())

	pageBytes, err := io.ReadAll(page)
	if err != nil {
		return 0, err
	}

	s.Add("page", &tengo.String{Value: string(pageBytes)})
	compiled, err := s.Compile()
	if err != nil {
		return 0, err
	}
	err = compiled.Run()
	if err != nil {
		return 0, err
	}
	return compiled.Get("result").Float(), nil

}
