package growthbook

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"go.klarlabs.de/rollops/pkg/flagconformance"
	"go.klarlabs.de/rollops/pkg/plugin"
)

func fakeGrowthBook(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestConformance(t *testing.T) {
	flagconformance.Run(t, func() (plugin.FlagProvider, error) {
		srv := fakeGrowthBook(t)
		return Provider{BaseURL: srv.URL, Token: "sec", HTTP: srv.Client()}, nil
	}, plugin.FlagChange{Flag: "checkout", Environment: "production"})
}
