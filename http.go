package sumoll

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

type httpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type (
	// HTTPSourceClient send to Resource HTTP
	HTTPSourceClient struct {
		url       *url.URL
		client    httpClient
		UserAgent string
		headers   *http.Header
	}

	// HTTPSourceClientOptFunc set some values to client
	HTTPSourceClientOptFunc func(*HTTPSourceClient) error
)

// SetXSumoCategoryHeader set X-Sumo-Category's value to header
func SetXSumoCategoryHeader(category string) HTTPSourceClientOptFunc {
	return func(h *HTTPSourceClient) error {
		if category != "" {
			h.headers.Add("X-Sumo-Category", category)
		}
		return nil
	}
}

// SetXSumoNameHeader set X-Sumo-Name's value to header
func SetXSumoNameHeader(name string) HTTPSourceClientOptFunc {
	return func(h *HTTPSourceClient) error {
		if name != "" {
			h.headers.Add("X-Sumo-Name", name)
		}
		return nil
	}
}

// SetXSumoHostHeader set X-Sumo-Host's value to header
func SetXSumoHostHeader(host string) HTTPSourceClientOptFunc {
	return func(h *HTTPSourceClient) error {
		if host != "" {
			h.headers.Add("X-Sumo-Host", host)
		}
		return nil
	}
}

// NewHTTPSourceClient create HTTPSourceClient object
func NewHTTPSourceClient(url *url.URL, opts ...HTTPSourceClientOptFunc) (*HTTPSourceClient, error) {
	c := &HTTPSourceClient{
		url:       url,
		client:    &http.Client{},
		UserAgent: UserAgent(),
		headers:   &http.Header{},
	}

	for _, opts := range opts {
		if err := opts(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Send object to sumologic
func (h *HTTPSourceClient) Send(body io.Reader) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	req, err := h.newRequest(ctx, http.MethodPost, body)
	if err != nil {
		return err
	}

	if h.headers != nil {
		mergeHeaders(&req.Header, h.headers)
	}

	res, err := h.client.Do(req)
	if err != nil {
		return err
	}

	if !validResponseStatus(res.StatusCode) {
		return errors.New(fmt.Sprintf("Unexpected response code from Sumologic: %v", res.StatusCode))
	}

	return nil
}

func validResponseStatus(status int) bool {
	return status >= http.StatusOK && status < http.StatusMultipleChoices
}

func mergeHeaders(merged, input *http.Header) {
	if input == nil {
		return
	}
	for k, v := range *input {
		if len(v) > 0 {
			merged.Add(k, v[0])
		}
	}
}

func (h *HTTPSourceClient) newRequest(ctx context.Context, method string, body io.Reader) (*http.Request, error) {
	u := *h.url

	req, err := http.NewRequest(method, u.String(), body)
	if err != nil {
		return nil, err
	}

	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", h.UserAgent)

	return req, nil
}
