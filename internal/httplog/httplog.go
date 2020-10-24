package httplog

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"os"
)

type Transport struct{}

func (t *Transport) RoundTrip(req *http.Request) (*http.Response, error) {
	dump, _ := httputil.DumpRequest(req, true)

	fmt.Fprintf(os.Stderr, "%q\n", dump)

	return http.DefaultTransport.RoundTrip(req)
}
