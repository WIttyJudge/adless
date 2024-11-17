package http

import (
	"io"
	"net"
	"net/http"
	"time"
)

const Timeout = 10 * time.Second

type HTTP struct {
	client *http.Client
}

func New() *HTTP {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: (&net.Dialer{
			Timeout: Timeout,
		}).Dial,
		IdleConnTimeout:     Timeout,
		TLSHandshakeTimeout: Timeout,
		MaxConnsPerHost:     10,
		MaxIdleConns:        10,
		MaxIdleConnsPerHost: 10,
	}

	client := &http.Client{
		Timeout:   Timeout,
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
