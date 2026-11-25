package handler

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/torbenish/clima-cep/internal/cep"
	"github.com/torbenish/clima-cep/internal/temperature"
)

type CEPLookup interface {
	Lookup(ctx context.Context, cep string) (string, error)
}

type WeatherProvider interface {
	CurrentTempC(ctx context.Context, city string) (float64, error)
}

type Handler struct {
	CEP     CEPLookup
	Weather WeatherProvider
}

type Response struct {
	TempC float64 `json:"temp_C"`
	TempF float64 `json:"temp_F"`
	TempK float64 `json:"temp_K"`
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	zip := r.PathValue("cep")

	if !cep.IsValid(zip) {
		http.Error(w, "invalid zipcode", http.StatusUnprocessableEntity)
		return
	}

	city, err := h.CEP.Lookup(r.Context(), zip)
	if errors.Is(err, cep.ErrNotFound) {
		http.Error(w, "can not find zipcode", http.StatusNotFound)
		return
	}
	if err != nil {
		log.Printf("cep lookup: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	c, err := h.Weather.CurrentTempC(r.Context(), city)
	if err != nil {
		log.Printf("weather: %v", err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		TempC: c,
		TempF: temperature.CelsiusToFahrenheit(c),
		TempK: temperature.CelsiusToKelvin(c),
	})
}