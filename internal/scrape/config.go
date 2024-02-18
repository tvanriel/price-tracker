package scrape

type Target struct {
	Cookies     map[string]string
	Headers     map[string][]string
	Page        string
        HTTPVersion string `mapstructure: "httpversion"`
	Compression bool

	Script string
	Labels map[string]string
	Name   string
}
type Config struct {
	Targets    []Target
	UserAgents []string
}
