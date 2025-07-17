package webhandler

import (
	"fmt"
	"net/http"

	"github.com/leonlonsdale/chirpy/internal/config"
)

func RegisterMetricsHandler(mux *http.ServeMux, cfg *config.ApiConfig) {
	mux.HandleFunc("GET /admin/metrics", MetricsHandler(cfg))
}

func MetricsHandler(cfg *config.ApiConfig) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		resString := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
			cfg.FileserverHits.Load())

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, resString)
	}
}
