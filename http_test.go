package sumoll

import (
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"
)

const os_env_http_source = "SUMOLL_TEST_HTTP_SOURCE_URL"

var (
	httpSourceURL = os.Getenv(os_env_http_source)
)

func TestHTTPSourceClientSendIntegration(t *testing.T) {
	if httpSourceURL == "" {
		t.Skip(os_env_http_source, "is not set. Value:", httpSourceURL)
	}
	t.Log(os_env_http_source, "is set. Executing integration test")

	table := []struct {
		body string
	}{
		{"hogehoge"},
	}

	for _, row := range table {
		u, err := url.Parse(httpSourceURL)
		if err != nil {
			t.Fatalf("url.Parse(%s) got error: %s", httpSourceURL, err)
		}
		c, err := NewHTTPSourceClient(u)
		if err != nil {
			t.Fatalf("NewHTTPSourceClient got error")
		}

		err = c.Send(strings.NewReader(row.body))
		if err != nil {
			t.Error("Error response when sending", row.body, "to", u, ". Err:", err)
		}
	}
}

type httpClientMock struct {
}

type nopCloser struct {
	io.Reader
}

func (nopCloser) Close() error { return nil }

const urlForUnitTests = "https://localEndpoint/receiver/v1/http/myUniqueID"

func (hc httpClientMock) Do(req *http.Request) (*http.Response, error) {
	operation := req.Method
	requestSize := req.ContentLength
	log.Printf("Mock intercepted %s request against %s with a body of %d bytes", operation, req.Host, requestSize)
	responseToSend := http.StatusBadRequest
	switch operation {
	case "POST":
		expectedURL, _ := url.Parse(urlForUnitTests)
		if req.URL != nil && *req.URL == *expectedURL {
			responseToSend = http.StatusOK
		}
	}
	return &http.Response{
		Body:       nopCloser{},
		StatusCode: responseToSend,
		Status:     http.StatusText(responseToSend),
	}, nil
}

func TestHTTPSourceClientSendUnit(t *testing.T) {
	table := []struct {
		category, hostname, sourcename string
	}{
		{"", "", ""},
	}

	for _, values := range table {
		localUrl, _ := url.Parse(urlForUnitTests)
		c, err := NewHTTPSourceClient(localUrl,
			SetXSumoCategoryHeader(values.category),
			SetXSumoHostHeader(values.hostname),
			SetXSumoNameHeader(values.sourcename),
		)
		if err != nil {
			t.Fatalf("NewHTTPSourceClient got error: %s", err)
		}

		c.client = &httpClientMock{}
		err = c.Send(strings.NewReader("hogehoge"))
		if err != nil {
			t.Error("Error response when sending payload to", localUrl,
				"with", values.category, values.hostname, values.sourcename,
				". Err:", err)
		}

	}
}
