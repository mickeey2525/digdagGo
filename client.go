package digdagGo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"runtime"
)

type Client struct {
	BaseURL    *url.URL
	HTTPClient *http.Client
	Token      string
	Logger     *log.Logger
}

var (
	version   = "0.0.1"
	userAgent = fmt.Sprintf("XXXGoClient/%s (%s)", version, runtime.Version())
)

func New(rawBaseURL, token string, logger *log.Logger) (*Client, error) {
	baseURL, err := url.Parse(rawBaseURL)
	if err != nil {
		return nil, err
	}
	if logger == nil {
		logger = log.New(os.Stderr, "[LOG]", log.LstdFlags)
	}

	return &Client{
		BaseURL:    baseURL,
		HTTPClient: http.DefaultClient,
		Token:      token,
		Logger:     logger,
	}, nil
}

func (c *Client) newRequest(ctx context.Context, method, spath string, params map[string]string, body io.Reader) (*http.Request, error) {
	reqURL := *c.BaseURL
	reqURL.Path = path.Join(reqURL.Path, spath)
	q := reqURL.Query()

	for k, v := range params {
		q.Add(k, v)
	}

	reqURL.RawQuery = q.Encode()
	switch method {
	case "GET":
		req, err := http.NewRequest(http.MethodGet, reqURL.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("TD1 %s", c.Token))
		req.Header.Add("User-Agent", userAgent)
		req = req.WithContext(ctx)
		return req, nil
	case "PUT":
		req, err := http.NewRequest(http.MethodPut, reqURL.String(), body)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("TD1 %s", c.Token))
		req.Header.Add("Content-Type", "application/gzip")
		req = req.WithContext(ctx)
		return req, nil
	case "DELETE":
		req, err := http.NewRequest(http.MethodDelete, reqURL.String(), nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("Authorization", fmt.Sprintf("TD1 %s", c.Token))
		req.Header.Add("Content-Type", "application/gzip")
		req = req.WithContext(ctx)
		return req, nil
	default:
		return nil, errors.New("you must specify method")

	}

}

func (c *Client) decodeBody(resp *http.Response, out interface{}) error {
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	return decoder.Decode(out)
}
