package api

import (
	"encoding/json"
	"errors"
	"log"
	"math"
	"net/http"
	"strconv"

	"github.com/Dhia-sel/convertix/internal/converter"
)

type convertResponse struct {
	Value  float64 `json:"value"`
	From   string  `json:"from"`
	To     string  `json:"to"`
	Result float64 `json:"result"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func NewHandler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", handleHealth)
	mux.HandleFunc("GET /convert", handleConvert)
	return mux
}

func handleHealth(w http.ResponseWriter, _ *http.Request) {
	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}

func handleConvert(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	valueStr := q.Get("value")
	from := q.Get("from")
	to := q.Get("to")

	if valueStr == "" || from == "" || to == "" {
		writeError(w, http.StatusBadRequest, "paramètres requis : value, from, to")
		return
	}

	value, err := strconv.ParseFloat(valueStr, 64)
	if err != nil || math.IsNaN(value) || math.IsInf(value, 0) {
		writeError(w, http.StatusBadRequest, "value doit être un nombre fini")
		return
	}

	result, err := converter.Convert(value, from, to)
	if err != nil {
		switch {
		case errors.Is(err, converter.ErrUnknownUnit):
			writeError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, converter.ErrIncompatibleUnits):
			writeError(w, http.StatusUnprocessableEntity, err.Error())
		default:
			log.Printf("convert: %v", err)
			writeError(w, http.StatusInternalServerError, "erreur interne")
		}
		return
	}

	if math.IsNaN(result) || math.IsInf(result, 0) {
		writeError(w, http.StatusUnprocessableEntity, "résultat hors limites")
		return
	}

	writeJSON(w, http.StatusOK, convertResponse{
		Value:  value,
		From:   from,
		To:     to,
		Result: result,
	})
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, errorResponse{Error: msg})
}

func writeJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	body, err := json.Marshal(v)
	if err != nil {
		log.Printf("writeJSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(`{"error":"erreur interne"}`))
		return
	}

	w.WriteHeader(status)
	_, _ = w.Write(body)
}
