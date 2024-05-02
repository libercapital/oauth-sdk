package oauthsdk

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"net/http"
	"sync"
	"time"

	"gitlab.com/bavatech/architecture/software/libs/go-modules/bavalogs.git"
	httptrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/net/http"
)

type client struct {
	Config Config
	Mutex  *sync.Mutex
	JWT    *JWT
	http   requestClient
}

type Client interface {
	GetAccessToken(ctx context.Context) (*JWT, error)
}

func applyCertificate(data CertData) (cfg *tls.Config, err error) {
	caKey, err := base64.StdEncoding.DecodeString(data.ClientCrtKey)
	if err != nil {
		return
	}

	caCert, err := base64.StdEncoding.DecodeString(data.ClientCrt)
	if err != nil {
		return
	}

	caCertPool := x509.NewCertPool()
	caCertPool.AppendCertsFromPEM(caCert)

	cert, err := tls.X509KeyPair(caCert, caKey)
	if err != nil {
		return
	}

	return &tls.Config{
		RootCAs:            caCertPool,
		Certificates:       []tls.Certificate{cert},
		InsecureSkipVerify: true,
	}, nil
}

func NewClient(config Config) (cli Client, err error) {
	tls := new(tls.Config)

	if config.CertData != nil {
		tls, err = applyCertificate(*config.CertData)
		if err != nil {
			return
		}
	}

	// Json is default
	if config.ContentType == "" {
		config.ContentType = ContentTypeApplicationJson
	}

	redactKeys := append(bavalogs.DefaultKeys, config.RedactKeys...)
	return &client{
		Config: config,
		Mutex:  &sync.Mutex{},
		JWT:    nil,
		http: httptrace.WrapClient(
			&http.Client{
				Transport: bavalogs.HttpClient{
					Proxied: &http.Transport{
						TLSClientConfig: tls,
					},
					RedactedKeys: redactKeys,
				},
				Timeout: 60 * time.Second,
			},
			httptrace.RTWithResourceNamer(func(req *http.Request) string {
				return "http.oauthsdk"
			}),
		),
	}, nil
}
