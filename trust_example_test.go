package trust_test

import (
	"crypto/tls"
	"net/http"

	"github.com/gomicro/trust"
)

var client *http.Client

func ExampleNew() {
	pool := trust.New()

	certs, err := pool.CACerts()
	if err != nil {
		panic(err)
	}

	client = &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{RootCAs: certs},
		},
	}
}
