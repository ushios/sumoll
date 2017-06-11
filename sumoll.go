package sumoll

import "fmt"

const (
	// Version string
	Version = "0.1"
)

// UserAgent string
func UserAgent() string {
	return fmt.Sprintf("Sumoll client/%s", Version)
}
