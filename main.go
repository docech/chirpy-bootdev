package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	router := chi.NewRouter()
	apiRouter := chi.NewRouter()
	adminRouter := chi.NewRouter()
	apiConfig := NewApiConfig()

	fsHandler := apiConfig.middlewareMetricsInc(http.StripPrefix("/app", http.FileServer(http.Dir("."))))
	
	adminRouter.Get("/metrics", apiConfig.metricsHandler)
	apiRouter.HandleFunc("/reset", apiConfig.resetHandler)
	apiRouter.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	router.Handle("/app/*", fsHandler)
	router.Handle("/app", fsHandler)
	router.Mount("/api", apiRouter)
	router.Mount("/admin", adminRouter)

	server := http.Server{
		Addr:    ":8080",
		Handler: corsMiddleware(router),
	}

	fmt.Println("Starting server at", server.Addr)
	server.ListenAndServe()
}

func corsMiddleware(next http.Handler) http.Handler {
	fmt.Println("CORS middleware")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PATCH, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}