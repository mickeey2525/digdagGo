package digdaggo

import (
	"archive/tar"
	"compress/gzip"
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
	"path/filepath"
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

func (c *Client) unarchive(dst string, r io.Reader) error {
	err := os.Mkdir(dst, 0755)
	if err != nil {
		return err
	}
	gzr, err := gzip.NewReader(r)
	if err != nil {
		return err
	}
	defer gzr.Close()
	tr := tar.NewReader(gzr)
	for {
		header, err := tr.Next()

		switch {
		case err == io.EOF:
			return nil

		case err != nil:
			return err

		case header == nil:
			continue
		}

		target := filepath.Join(dst, header.Name)

		// check the file type
		switch header.Typeflag {

		// if its a dir and it doesn't exist create it
		case tar.TypeDir:
			if _, err := os.Stat(target); err != nil {
				if err := os.MkdirAll(target, 0755); err != nil {
					return err
				}
			}

		// if it's a file create it
		case tar.TypeReg:
			f, err := os.OpenFile(target, os.O_CREATE|os.O_RDWR, os.FileMode(header.Mode))
			if err != nil {
				return err
			}

			// copy over contents
			if _, err := io.Copy(f, tr); err != nil {
				return err
			}

			// manually close here after each file operation; defering would cause each file close
			// to wait until all operations have completed.
			f.Close()
		}
		return nil
	}
}
