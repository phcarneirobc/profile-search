package client

import (
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"time"

	"github.com/phcarneirobc/profile-search/internal/config"
)

func CreateClient(config *config.Config) *http.Client {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 100,
	}

	if len(config.ProxyList) > 0 {
		proxyURL, err := url.Parse(config.ProxyList[rand.Intn(len(config.ProxyList))])
		if err == nil {
			transport.Proxy = http.ProxyURL(proxyURL)
		}
	}

	return &http.Client{
		Transport: transport,
		Timeout:   time.Duration(config.TimeoutSeconds) * time.Second,
	}
}

func ReadResponseBody(resp *http.Response) (string, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(body), nil
}
