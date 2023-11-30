package main

import (
	"fmt"
	"net/http"
)

type apiConfig struct {
	fileServerHits int
}

func NewApiConfig() *apiConfig {
	return &apiConfig{
		fileServerHits: 0,
	}
}

func (a *apiConfig) middlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		a.fileServerHits++

		next.ServeHTTP(w, r)
	})
}

const metricsTemplate = `
	<html>
	<body>
		<h1>Welcome, Chirpy Admin</h1>
		<p>Chirpy has been visited %d times!</p>
	</body>

	</html>
`


func (a* apiConfig) metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf(metricsTemplate, a.fileServerHits)))
}

func (a* apiConfig) resetHandler(w http.ResponseWriter, r *http.Request) {
	a.fileServerHits = 0
	w.WriteHeader(http.StatusOK)
}