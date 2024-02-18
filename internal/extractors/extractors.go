package extractors

import (
	"io"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

type extractorFunc func(string, string) (float64, error)

func innerTextExtract(selector string, page io.Reader) (float64, error) {
	doc, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		return 0, err
	}
	var price float64
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		pricetag := s.Text()
		f, err := cleanupPrice(pricetag)
		if err != nil {
			return
		}
		price = f
	})
	return price, nil
}
func attributeExtract(selector string, attribute string, page io.Reader) (float64, error) {

	doc, err := goquery.NewDocumentFromReader(page)
	if err != nil {
		return 0, err
	}
	var price float64
	doc.Find(selector).Each(func(i int, s *goquery.Selection) {
		pricetag, exists := s.Attr(attribute)
		if !exists {
			return
		}
		f, err := cleanupPrice(pricetag)
		if err != nil {
			return
		}
		price = f
	})
	return price, nil

}

func cleanupPrice(pricetag string) (float64, error) {
	var err error

	pricetag = strings.ReplaceAll(pricetag, ",", ".")
	pricetag = strings.TrimPrefix(pricetag, "$")
	pricetag = strings.TrimPrefix(pricetag, "€")
	pricetag = strings.TrimSpace(pricetag)
	var f float64
	if f, err = strconv.ParseFloat(pricetag, 64); err != nil {
		return 0.0, err
	}
	return f, nil
}
