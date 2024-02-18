package curlhttp_test

import (
	"testing"

	"github.com/tvanriel/price-tracker/internal/curlhttp"
	"gotest.tools/v3/assert"
)

var goldenSample = `* Host goo.gl:443 was resolved.
* IPv6: 2a00:1450:4001:827::200e
* IPv4: 172.217.16.142
*   Trying [2a00:1450:4001:827::200e]:443...
* Connected to goo.gl (2a00:1450:4001:827::200e) port 443
* ALPN: curl offers http/1.1
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
*  CAfile: /etc/ssl/certs/ca-certificates.crt
*  CApath: none
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, CERT verify (15):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS change cipher, Change cipher spec (1):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_256_GCM_SHA384 / x25519 / id-ecPublicKey
* ALPN: server accepted http/1.1
* Server certificate:
*  subject: CN=*.google.com
*  start date: Jan 29 08:04:47 2024 GMT
*  expire date: Apr 22 08:04:46 2024 GMT
*  subjectAltName: host "goo.gl" matched cert's "goo.gl"
*  issuer: C=US; O=Google Trust Services LLC; CN=GTS CA 1C3
*  SSL certificate verify ok.
*   Certificate level 0: Public key type EC/prime256v1 (256/128 Bits/secBits), signed using sha256WithRSAEncryption
*   Certificate level 1: Public key type RSA (2048/112 Bits/secBits), signed using sha256WithRSAEncryption
*   Certificate level 2: Public key type RSA (4096/152 Bits/secBits), signed using sha384WithRSAEncryption
* using HTTP/1.x
> GET / HTTP/1.1
> Host: goo.gl
> User-Agent: curl/8.6.0
> Accept: */*
>
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* TLSv1.3 (IN), TLS handshake, Newsession Ticket (4):
* old SSL session ID is stale, removing
< HTTP/1.1 301 Moved Permanently
< Content-Type: application/binary
< Vary: Sec-Fetch-Dest, Sec-Fetch-Mode, Sec-Fetch-Site
< Cache-Control: no-cache, no-store, max-age=0, must-revalidate
< Pragma: no-cache
< Expires: Mon, 01 Jan 1990 00:00:00 GMT
< Date: Sat, 17 Feb 2024 07:56:27 GMT
< Location: https://developers.googleblog.com/2018/03/transitioning-google-url-shortener.html
< Strict-Transport-Security: max-age=31536000
< Cross-Origin-Opener-Policy: unsafe-none
< Cross-Origin-Resource-Policy: same-site
< Accept-CH: Sec-CH-UA-Arch, Sec-CH-UA-Bitness, Sec-CH-UA-Full-Version, Sec-CH-UA-Full-Version-List, Sec-CH-UA-Model, Sec-CH-UA-WoW64, Sec-CH-UA-Form-Factor, Sec-CH-UA-Platform, Sec-CH-UA-Platform-Version
< Permissions-Policy: ch-ua-arch=*, ch-ua-bitness=*, ch-ua-full-version=*, ch-ua-full-version-list=*, ch-ua-model=*, ch-ua-wow64=*, ch-ua-form-factor=*, ch-ua-platform=*, ch-ua-platform-version=*
< Content-Security-Policy: require-trusted-types-for 'script';report-uri /_/DurableDeepLinkUi/cspreport
< Content-Security-Policy: script-src 'nonce-yJ70rmbupwCQd8DxJ_NX8A' 'unsafe-inline';object-src 'none';base-uri 'self';report-uri /_/DurableDeepLinkUi/cspreport;worker-src 'self'
< Server: ESF
< Content-Length: 0
< X-XSS-Protection: 0
< X-Frame-Options: SAMEORIGIN
< X-Content-Type-Options: nosniff
< Alt-Svc: h3=":443"; ma=2592000,h3-29=":443"; ma=2592000
<
* Connection #0 to host goo.gl left intact
`

func TestCode(t *testing.T) {
	code, err := curlhttp.CurlHttpStatusCodeParser(goldenSample)
	assert.NilError(t, err)
	assert.Equal(t, "301", code)
}

func TestHeaders(t *testing.T) {
	headers := curlhttp.ParseHTTPHeaders(goldenSample)
	assert.Equal(t, len(headers), 19)

	assert.Equal(t, len(headers["Vary"]), 1)
	assert.Equal(t, headers["Vary"][0], "Sec-Fetch-Dest, Sec-Fetch-Mode, Sec-Fetch-Site")

	assert.Equal(t, len(headers["Server"]), 1)
	assert.Equal(t, headers["Server"][0], "ESF")

	assert.Equal(t, len(headers["Content-Security-Policy"]), 7)
}
