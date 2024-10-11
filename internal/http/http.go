package http

import (
	"barrier/internal/config"
	"crypto/tls"
	"io"
	"net"
	"net/http"
)

type HTTP struct {
	client *http.Client
}

func New(config *config.Config) *HTTP {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout: config.HTTP.Timeout,
		}).Dial,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		IdleConnTimeout:     config.HTTP.Timeout,
		TLSHandshakeTimeout: config.HTTP.Timeout,
		MaxConnsPerHost:     10,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}

	client := &http.Client{
		Timeout:   config.HTTP.Timeout,
		Transport: transport,
	}

	return &HTTP{
		client: client,
	}
}

func (h *HTTP) Get(url string) (string, error) {
	resp, err := h.client.Get(url)
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(body), nil
}
