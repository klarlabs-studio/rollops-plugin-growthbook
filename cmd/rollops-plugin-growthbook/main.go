// Command rollops-plugin-growthbook is a Rollops feature-flag provider plugin
// backed by GrowthBook. Build it, pin its sha256, and point a rollout's
// featureFlags.plugin at the binary.
package main

import (
	"fmt"
	"os"

	growthbook "github.com/klarlabs-studio/rollops-plugin-growthbook"
	"go.klarlabs.de/rollops/pkg/plugin"
)

// version is overwritten at build time via -ldflags.
var version = "dev"

func main() {
	safety := plugin.Safety{
		NetworkHosts: []string{"api.growthbook.io:443"},
		EnvVars:      []string{"GROWTHBOOK_API_URL", "GROWTHBOOK_TOKEN"},
		RiskClass:    plugin.RiskActive,
	}
	if err := plugin.ServeFlagProvider("klarlabs/growthbook", version, growthbook.FromEnv(), safety); err != nil {
		fmt.Fprintln(os.Stderr, "rollops-plugin-growthbook:", err)
		os.Exit(1)
	}
}
