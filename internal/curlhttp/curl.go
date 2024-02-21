package curlhttp

import (
	"bytes"
	"context"
	"io"
	"os/exec"
	"strings"

	"go.uber.org/zap"
)

type NewCurlHttpParams struct {
	Config Config
	Log    *zap.Logger
}

type CurlHttp struct {
	Config Config
	Log    *zap.Logger
}

func NewCurlHttp(p NewCurlHttpParams) *CurlHttp {
	return &CurlHttp{
		Config: p.Config,
		Log:    p.Log.Named("curlhttp"),
	}
}

func (c *CurlHttp) Do(ctx context.Context, req Request) (Response, error) {

	args := []string{
		"-A", req.UserAgent(),
		"-vvv",
	}
	if req.Method() != "GET" {
		args = append(args, "-X", req.Method())
	}

	if req.Compressed() {
		args = append(args, "--compressed")
	}
	args = append(args, strings.Join([]string{"--http", req.HTTPVersion()}, ""))

	headers := req.Headers()
	for k := range headers {
                args = append(args, "-H", strings.Join([]string{k, ": ", strings.Join(headers[k], ";")}, ""))

	}

	args = append(args, req.Url())

	buf := bytes.NewBuffer([]byte{})
	errbuf := bytes.NewBuffer([]byte{})
	cmd := exec.CommandContext(ctx, c.Config.Binary, args...)

	cmd.Stdout = buf
	cmd.Stderr = errbuf
	cmd.Stdin = req.Body()
	c.Log.Info("GET", zap.Strings("args", args))

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	code, err := CurlHttpStatusCodeParser(errbuf.String())
	if err != nil {
		return nil, err
	}
	responseHeaders := ParseHTTPHeaders(errbuf.String())

	resp := &CurlResponse{
		body:     buf,
		httpCode: code,
		headers:  responseHeaders,
	}
	return resp, nil

}

type Request interface {
	HTTPVersion() string
	Method() string
	Url() string
	Headers() map[string][]string
	Cookies() map[string]string
	UserAgent() string
	Compressed() bool
	Body() io.Reader
}

type CurlRequest struct {
	httpVersion string
	method      string
	url         string
	cookies     map[string]string
	userAgent   string
	headers     map[string][]string
	compressed  bool
	body        io.Reader
}

func (c *CurlRequest) HTTPVersion() string          { return c.httpVersion }
func (c *CurlRequest) Method() string               { return c.method }
func (c *CurlRequest) Url() string                  { return c.url }
func (c *CurlRequest) Cookies() map[string]string   { return c.cookies }
func (c *CurlRequest) UserAgent() string            { return c.userAgent }
func (c *CurlRequest) Headers() map[string][]string { return c.headers }
func (c *CurlRequest) Body() io.Reader              { return c.body }
func (c *CurlRequest) Compressed() bool             { return c.compressed }

func (c *CurlRequest) EnableCompression()  { c.compressed = true }
func (c *CurlRequest) DisableCompression() { c.compressed = false }
func (c *CurlRequest) HTTP3()              { c.httpVersion = "3" }
func (c *CurlRequest) HTTP2()              { c.httpVersion = "2" }
func (c *CurlRequest) HTTP1()              { c.httpVersion = "1.1" }
func (c *CurlRequest) SetURL(url string)   { c.url = url }
func (c *CurlRequest) AddHeader(name string, value []string) {
	if c.headers == nil {
		c.headers = make(map[string][]string)
	}
	c.headers[name] = value
}
func (c *CurlRequest) SetMethod(method string)              { c.method = method }
func (c *CurlRequest) SetCookies(cookies map[string]string) { c.cookies = cookies }
func (c *CurlRequest) SetUserAgent(useragent string)        { c.userAgent = useragent }
func (c *CurlRequest) SetBody(reader io.Reader)             { c.body = reader }

// Type assert
var _ Request = &CurlRequest{}

type Response interface {
	Body() io.Reader
	HttpCode() string
	Headers() map[string][]string
}

type CurlResponse struct {
	body     io.Reader
	httpCode string
	headers  map[string][]string
}

func (c *CurlResponse) Body() io.Reader              { return c.body }
func (c *CurlResponse) HttpCode() string             { return c.httpCode }
func (c *CurlResponse) Headers() map[string][]string { return c.headers }

func joinKeyValuePairsToString(inputMap map[string]string) string {
	resultSlice := make([]string, len(inputMap))
	index := 0
	for key, value := range inputMap {
		resultSlice[index] = strings.Join([]string{key, "=", value}, "")
		index++
	}
	return strings.Join(resultSlice, ";")
}
