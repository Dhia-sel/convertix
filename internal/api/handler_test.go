package api

import (
	"encoding/json"
	"math"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHealthz(t *testing.T) {
	handler := NewHandler()
	req := httptest.NewRequest(http.MethodGet, "/healthz", nil)
	rec := httptest.NewRecorder()

	handler.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Errorf("code = %d, attendu %d", rec.Code, http.StatusOK)
	}
	if got := rec.Header().Get("Content-Type"); got != "application/json" {
		t.Errorf("Content-Type = %q, attendu %q", got, "application/json")
	}
	if got := rec.Body.String(); got != `{"status":"ok"}` {
		t.Errorf("corps = %q, attendu %q", got, `{"status":"ok"}`)
	}
}

func TestConvertEndpoint(t *testing.T) {
	handler := NewHandler()

	tests := []struct {
		name       string
		method     string
		query      url.Values
		wantStatus int
		wantResult float64
	}{
		{name: "1 mètre en centimètres", method: http.MethodGet, query: url.Values{"value": {"1"}, "from": {"mètre"}, "to": {"centimètre"}}, wantStatus: http.StatusOK, wantResult: 100},
		{name: "1 mètre en pouces", method: http.MethodGet, query: url.Values{"value": {"1"}, "from": {"mètre"}, "to": {"pouce"}}, wantStatus: http.StatusOK, wantResult: 39.37007874015748},
		{name: "même unité", method: http.MethodGet, query: url.Values{"value": {"5"}, "from": {"mètre"}, "to": {"mètre"}}, wantStatus: http.StatusOK, wantResult: 5},
		{name: "paramètre to manquant", method: http.MethodGet, query: url.Values{"value": {"1"}, "from": {"mètre"}}, wantStatus: http.StatusBadRequest},
		{name: "value non numérique", method: http.MethodGet, query: url.Values{"value": {"abc"}, "from": {"mètre"}, "to": {"pouce"}}, wantStatus: http.StatusBadRequest},
		{name: "value vaut NaN", method: http.MethodGet, query: url.Values{"value": {"NaN"}, "from": {"mètre"}, "to": {"pouce"}}, wantStatus: http.StatusBadRequest},
		{name: "depuis unité inconnue", method: http.MethodGet, query: url.Values{"value": {"1"}, "from": {"unité_inconnue"}, "to": {"mètre"}}, wantStatus: http.StatusBadRequest},
		{name: "vers unité inconnue", method: http.MethodGet, query: url.Values{"value": {"1"}, "from": {"mètre"}, "to": {"unité_inconnue"}}, wantStatus: http.StatusBadRequest},
		{name: "mètre vers kilogramme (incompatible)", method: http.MethodGet, query: url.Values{"value": {"1"}, "from": {"mètre"}, "to": {"kilogramme"}}, wantStatus: http.StatusUnprocessableEntity},
		{name: "résultat hors limites", method: http.MethodGet, query: url.Values{"value": {"1e308"}, "from": {"tonne longue"}, "to": {"grain"}}, wantStatus: http.StatusUnprocessableEntity},
		{name: "mauvaise méthode", method: http.MethodPost, query: url.Values{"value": {"1"}, "from": {"mètre"}, "to": {"pouce"}}, wantStatus: http.StatusMethodNotAllowed},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			req := httptest.NewRequest(tc.method, "/convert?"+tc.query.Encode(), nil)
			rec := httptest.NewRecorder()

			handler.ServeHTTP(rec, req)

			if rec.Code != tc.wantStatus {
				t.Fatalf("code = %d, attendu %d", rec.Code, tc.wantStatus)
			}

			switch tc.wantStatus {
			case http.StatusOK:
				var resp convertResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("JSON illisible : %v", err)
				}
				if math.Abs(resp.Result-tc.wantResult) > 1e-9 {
					t.Errorf("résultat = %v, attendu %v", resp.Result, tc.wantResult)
				}
			case http.StatusBadRequest, http.StatusUnprocessableEntity:
				var resp errorResponse
				if err := json.NewDecoder(rec.Body).Decode(&resp); err != nil {
					t.Fatalf("JSON d'erreur illisible : %v", err)
				}
				if resp.Error == "" {
					t.Error("le message d'erreur est vide")
				}
			}
		})
	}
}
