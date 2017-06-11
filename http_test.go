package sumoll

import (
	"net/url"
	"os"
	"strings"
	"testing"
)

var (
	httpSourceURL = os.Getenv("SUMOLL_TEST_HTTP_SOURCE_URL")
)

func httpSourceTestAvailable() bool {
	if httpSourceURL == "" {
		return false
	}

	return true
}

func TestHTTPSourceClientSend(t *testing.T) {
	if !httpSourceTestAvailable() {
		t.Skip("SUMOLL_TEST_HTTP_SOURCE_URL is not set")
	}

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
		c := NewHTTPSourceClient(u)

		c.Send(strings.NewReader(row.body))
	}
}
