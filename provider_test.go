package growthbook

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"go.klarlabs.de/rollops/pkg/plugin"
)

func TestApplyFlag_PatchesEnvironmentRule(t *testing.T) {
	var path string
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path = r.URL.Path
		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &body)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	p := Provider{BaseURL: srv.URL, Token: "sec", HTTP: srv.Client()}
	if err := p.ApplyFlag(context.Background(), plugin.FlagChange{Flag: "checkout", Environment: "production", Percentage: 40}); err != nil {
		t.Fatalf("ApplyFlag: %v", err)
	}
	if !strings.HasSuffix(path, "/api/v1/features/checkout") {
		t.Errorf("wrong path: %s", path)
	}
	envs, _ := body["environments"].(map[string]any)
	prod, _ := envs["production"].(map[string]any)
	if prod["enabled"] != true {
		t.Errorf("enabled = %v, want true", prod["enabled"])
	}
	rules, _ := prod["rules"].([]any)
	r0, _ := rules[0].(map[string]any)
	if r0["type"] != "rollout" || r0["coverage"].(float64) != 0.4 {
		t.Errorf("rule = %v, want rollout coverage 0.4", r0)
	}
}

func TestApplyFlag_DisabledSendsFalse(t *testing.T) {
	var body map[string]any
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		json.Unmarshal(b, &body)
		w.WriteHeader(200)
	}))
	defer srv.Close()

	p := Provider{BaseURL: srv.URL, Token: "sec", HTTP: srv.Client()}
	if err := p.ApplyFlag(context.Background(), plugin.FlagChange{Flag: "f", Environment: "staging", Percentage: 0, Disabled: true}); err != nil {
		t.Fatalf("ApplyFlag: %v", err)
	}
	envs, _ := body["environments"].(map[string]any)
	st, _ := envs["staging"].(map[string]any)
	if st["enabled"] != false {
		t.Errorf("enabled = %v, want false", st["enabled"])
	}
}

func TestApplyFlag_ServerErrorPropagates(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(403)
	}))
	defer srv.Close()
	p := Provider{BaseURL: srv.URL, Token: "sec", HTTP: srv.Client()}
	if err := p.ApplyFlag(context.Background(), plugin.FlagChange{Flag: "f", Environment: "p"}); err == nil {
		t.Fatal("403 must error")
	}
}

func TestApplyFlag_RequiresToken(t *testing.T) {
	p := Provider{BaseURL: "http://x"}
	if err := p.ApplyFlag(context.Background(), plugin.FlagChange{Flag: "f", Environment: "p"}); err == nil {
		t.Fatal("missing token must error")
	}
}
