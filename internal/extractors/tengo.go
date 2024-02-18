package extractors

import (
	"strings"

	tengo "github.com/d5/tengo/v2"
	"go.uber.org/zap"
)

func TengoInnerTextExtractor(selector string, log *zap.Logger) *tengo.UserFunction {

	log = log.With(zap.String("selector", selector))
	return &tengo.UserFunction{
		Name: "innerText",
		Value: func(args ...tengo.Object) (ret tengo.Object, err error) {
			if len(args) != 1 {
				return nil, tengo.ErrWrongNumArguments
			}
			s, ok := tengo.ToString(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "page",
					Expected: "string",
					Found:    args[0].TypeName(),
				}
			}

			price, err := innerTextExtract(selector, strings.NewReader(s))

			if err != nil {
				log.Info("errored while extracting innerHTML",
					zap.Error(err),
				)
				return nil, err
			}

			log.Info("extracted innerHTML",
				zap.Float64("price", price),
			)
			return &tengo.Float{Value: price}, nil
		}}

}
func TengoAttributeExtractor(selector string, attribute string, log *zap.Logger) *tengo.UserFunction {
	log = log.With(zap.String("attribute", attribute), zap.String("selector", selector))
	return &tengo.UserFunction{
		Name: "attr",
		Value: func(args ...tengo.Object) (tengo.Object, error) {
			if len(args) != 1 {
				return nil, tengo.ErrWrongNumArguments
			}
			s, ok := tengo.ToString(args[0])
			if !ok {
				return nil, tengo.ErrInvalidArgumentType{
					Name:     "page",
					Expected: "string",
					Found:    args[0].TypeName(),
				}
			}

			price, err := attributeExtract(selector, attribute, strings.NewReader(s))
			if err != nil {
				log.Info("errored while extracting attribute",
					zap.Error(err),
				)
				return nil, err
			}

			log.Info("extracted attribute",
				zap.Float64("price", price),
			)

			return &tengo.Float{Value: price}, nil
		}}

}
