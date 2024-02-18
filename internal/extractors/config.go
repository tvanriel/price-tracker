package extractors

type Config struct {
	Attribute map[string]struct {
		Selector  string
		Attribute string
	} `mapstructure:"attribute"`
	InnerText map[string]struct {
		Selector string
	} `mapstructure:"innerText"`
}
