package handler

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/torbenish/clima-cep/internal/cep"
)

type fakeCEP struct {
	city string
	err  error
}

func (f fakeCEP) Lookup(ctx context.Context, _ string) (string, error) { return f.city, f.err }

type fakeWeather struct{ temp float64 }

func (f fakeWeather) CurrentTempC(ctx context.Context, _ string) (float64, error) {
	return f.temp, nil
}

func do(h *Handler, path string) *httptest.ResponseRecorder {
	mux := http.NewServeMux()
	mux.Handle("GET /{cep}", h)
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestSuccess(t *testing.T) {
	h := &Handler{CEP: fakeCEP{city: "São Paulo"}, Weather: fakeWeather{temp: 28.5}}
	rec := do(h, "/01001000")

	if rec.Code != http.StatusOK {
		t.Fatalf("status %d", rec.Code)
	}
	body := rec.Body.String()
	for _, want := range []string{`"temp_C":28.5`, `"temp_F":83.3`, `"temp_K":301.5`} {
		if !strings.Contains(body, want) {
			t.Errorf("body %q missing %s", body, want)
		}
	}
}

func TestInvalidZipcode(t *testing.T) {
	h := &Handler{CEP: fakeCEP{}, Weather: fakeWeather{}}
	for _, path := range []string{"/123", "/1234567a", "/123456789"} {
		rec := do(h, path)
		if rec.Code != http.StatusUnprocessableEntity ||
			strings.TrimSpace(rec.Body.String()) != "invalid zipcode" {
			t.Errorf("%s: got %d %q", path, rec.Code, rec.Body.String())
		}
	}
}

func TestZipcodeNotFound(t *testing.T) {
	h := &Handler{CEP: fakeCEP{err: cep.ErrNotFound}, Weather: fakeWeather{}}
	rec := do(h, "/99999999")
	if rec.Code != http.StatusNotFound ||
		strings.TrimSpace(rec.Body.String()) != "can not find zipcode" {
		t.Errorf("got %d %q", rec.Code, rec.Body.String())
	}
}