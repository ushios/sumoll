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
	checkCategoryValue, checkHostnameValue, checkSourcenameValue string
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
		if req.URL == nil || *req.URL != *expectedURL {
			log.Println("Request against mock did not match the expected URL", expectedURL)
		} else if hc.checkCategoryValue != "" && !hasHeaderWithValue(req, "X-Sumo-Category", hc.checkCategoryValue) {
			log.Println("Request against mock did not have the expected category value", hc.checkCategoryValue)
		} else if hc.checkHostnameValue != "" && !hasHeaderWithValue(req, "X-Sumo-Host", hc.checkHostnameValue) {
			log.Println("Request against mock did not have the expected hostname value", hc.checkHostnameValue)
		} else if hc.checkSourcenameValue != "" && !hasHeaderWithValue(req, "X-Sumo-Name", hc.checkSourcenameValue) {
			log.Println("Request against mock did not have the expected sourcename value", hc.checkSourcenameValue)
		} else {
			log.Println("All mock tests were successful")
			responseToSend = http.StatusOK
		}
	}
	return &http.Response{
		Body:       nopCloser{},
		StatusCode: responseToSend,
		Status:     http.StatusText(responseToSend),
	}, nil
}
func hasHeaderWithValue(request *http.Request, header string, expectedValue string) bool {
	return request.Header.Get(header) == expectedValue
}

func TestHTTPSourceClientSendUnit(t *testing.T) {
	table := []struct {
		category, hostname, sourcename string
	}{
		{"", "", ""},
		{"testCategory", "", ""},
		{"", "testHostname", ""},
		{"", "", "testSourcename"},
		{"testCategory", "testHostname", "testSourcename"},
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

		c.client = &httpClientMock{
			checkCategoryValue:   values.category,
			checkHostnameValue:   values.hostname,
			checkSourcenameValue: values.sourcename,
		}
		err = c.Send(strings.NewReader("hogehoge"))
		if err != nil {
			t.Error("Error response when sending payload to", localUrl,
				"with", values.category, values.hostname, values.sourcename,
				". Err:", err)
		}

	}
}
