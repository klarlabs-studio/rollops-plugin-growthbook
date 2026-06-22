// Package growthbook is a Rollops feature-flag provider plugin backed by
// GrowthBook's REST API. It drives a feature's enabled state and a single
// rollout rule's coverage to match a rollout's progressive steps, so a
// GrowthBook flag tracks a Rollops canary in lockstep.
package growthbook

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.klarlabs.de/rollops/pkg/plugin"
)

// Provider talks to GrowthBook's REST API. BaseURL and Token come from the
// plugin's environment (see Config); Environment is supplied per call by Rollops
// as the GrowthBook environment id (e.g. "production").
type Provider struct {
	BaseURL string // e.g. https://api.growthbook.io
	Token   string // API secret (Authorization: Bearer <secret>)
	HTTP    *http.Client
}

func (p Provider) client() *http.Client {
	if p.HTTP != nil {
		return p.HTTP
	}
	return http.DefaultClient
}

// ApplyFlag PATCHes the feature so the given environment is enabled/disabled and
// carries one rollout rule whose coverage equals the percentage. GrowthBook
// merges the per-environment patch, leaving other environments untouched.
func (p Provider) ApplyFlag(ctx context.Context, c plugin.FlagChange) error {
	if p.Token == "" {
		return fmt.Errorf("growthbook: GROWTHBOOK_TOKEN is required")
	}
	rule := map[string]any{
		"type":          "rollout",
		"description":   "rollops canary",
		"value":         "true",
		"coverage":      float64(c.Percentage) / 100.0,
		"hashAttribute": "id",
		"enabled":       true,
	}
	body := map[string]any{
		"environments": map[string]any{
			c.Environment: map[string]any{
				"enabled": !c.Disabled,
				"rules":   []any{rule},
			},
		},
	}
	u := fmt.Sprintf("%s/api/v1/features/%s", p.BaseURL, url.PathEscape(c.Flag))
	if err := p.do(ctx, http.MethodPost, u, body); err != nil {
		return fmt.Errorf("growthbook: update feature %q: %w", c.Flag, err)
	}
	return nil
}

func (p Provider) do(ctx context.Context, method, u string, body any) error {
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, method, u, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+p.Token)
	req.Header.Set("Content-Type", "application/json")
	resp, err := p.client().Do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()
	if resp.StatusCode >= 300 {
		return fmt.Errorf("status %d", resp.StatusCode)
	}
	return nil
}
