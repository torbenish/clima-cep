package cep

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsValid(t *testing.T) {
	tests := map[string]bool{
		"01001000":  true,
		"1234567":   false,
		"123456789": false,
		"1234567a":  false,
		"12345-678": false,
		"":          false,
	}
	for in, want := range tests {
		if got := IsValid(in); got != want {
			t.Errorf("IsValid(%q) = %v, want %v", in, got, want)
		}
	}
}

func TestLookup(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/ws/01001000/json/":
			w.Write([]byte(`{"localidade":"São Paulo"}`))
		default:
			w.Write([]byte(`{"erro":"true"}`))
		}
	}))
	defer srv.Close()

	c := &Client{BaseURL: srv.URL, HTTP: srv.Client()}

	city, err := c.Lookup(context.Background(), "01001000")
	if err != nil || city != "São Paulo" {
		t.Fatalf("got %q, %v", city, err)
	}
	if _, err := c.Lookup(context.Background(), "99999999"); err != ErrNotFound {
		t.Fatalf("want ErrNotFound, got %v", err)
	}
}