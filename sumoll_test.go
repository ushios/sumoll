package sumoll

import "testing"

func TestUserAgent(t *testing.T) {
	ua := UserAgent()

	if ua == "" {
		t.Errorf("UserAgent is empty")
	}
}
