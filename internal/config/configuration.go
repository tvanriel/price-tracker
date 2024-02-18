package config

import (
	"strings"

	"github.com/spf13/viper"
	"github.com/tvanriel/cloudsdk/logging"
	"github.com/tvanriel/cloudsdk/prometheus"
	"github.com/tvanriel/price-tracker/internal/curlhttp"
	"github.com/tvanriel/price-tracker/internal/extractors"
	"github.com/tvanriel/price-tracker/internal/scrape"
	"go.uber.org/fx"
)

type AppConfig struct {
	Prometheus prometheus.Configuration `mapstructure:"prometheus"`
	Logging    logging.Configuration    `mapstructure:"log"`
	Scrape     scrape.Config            `mapstructure:"scrape"`
	Extractors extractors.Config        `mapstructure:"extractors"`
	Curl       curlhttp.Config          `mapstructure:"curl"`
}

func Configuration() fx.Option {
	return fx.Provide(
		ViperConfiguration,
		PrometheusConfiguration,
		ScrapeConfiguration,
		LoggingConfiguration,
		ExtractorConfiguration,
		CurlConfiguration,
	)
}

func ViperConfiguration() (AppConfig, error) {

	var config AppConfig
	v := viper.NewWithOptions(viper.KeyDelimiter("::"))
	v.AddConfigPath(".")
	v.AddConfigPath("/etc/scrape")
	v.SetConfigName("config")
	v.SetConfigType("yaml")
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	err := v.ReadInConfig()
	if err != nil {
		return AppConfig{}, err
	}
	err = v.Unmarshal(&config)
	if err != nil {
		return AppConfig{}, err
	}

	return config, err
}
func PrometheusConfiguration(c AppConfig) prometheus.Configuration {
	return c.Prometheus
}

func ScrapeConfiguration(c AppConfig) scrape.Config {
	return c.Scrape
}
func ExtractorConfiguration(c AppConfig) extractors.Config {
	return c.Extractors
}
func LoggingConfiguration(c AppConfig) logging.Configuration {
	return c.Logging
}
func CurlConfiguration(c AppConfig) curlhttp.Config {
	return c.Curl
}
