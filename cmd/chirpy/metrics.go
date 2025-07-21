package main

import (
	"fmt"
	"net/http"
)

func (app *application) MiddlewareMetricsInc(next http.Handler) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			app.config.FileserverHits.Add(1)
			next.ServeHTTP(w, r)
		})
}

func (app *application) MetricsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resString := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
			app.config.FileserverHits.Load())

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, resString)

	}
}
