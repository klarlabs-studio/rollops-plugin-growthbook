package growthbook

import "os"

// FromEnv builds a Provider from the plugin's environment. Secrets and endpoint
// come from the plugin process, never from the Rollops target spec (Rollops
// passes only the flag name, environment, and percentage).
//
//	GROWTHBOOK_API_URL  base URL (default https://api.growthbook.io)
//	GROWTHBOOK_TOKEN    API secret (required)
func FromEnv() Provider {
	base := os.Getenv("GROWTHBOOK_API_URL")
	if base == "" {
		base = "https://api.growthbook.io"
	}
	return Provider{
		BaseURL: base,
		Token:   os.Getenv("GROWTHBOOK_TOKEN"),
	}
}
