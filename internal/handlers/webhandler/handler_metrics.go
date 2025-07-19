package webhandler

import (
	"fmt"
	"net/http"
	"sync/atomic"

	"github.com/leonlonsdale/chirpy/internal/database"
)

func MetricsHandler(db database.Queries, fs *atomic.Int32) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		resString := fmt.Sprintf(`
		<html>
			<body>
				<h1>Welcome, Chirpy Admin</h1>
				<p>Chirpy has been visited %d times!</p>
			</body>
		</html>`,
			fs.Load())

		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		_, _ = fmt.Fprint(w, resString)
	}
}
