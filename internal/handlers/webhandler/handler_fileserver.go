package webhandler

import (
	"net/http"
)

func FileServerHandler(mux *http.ServeMux, middleware func(next http.Handler) http.Handler) {
	mux.Handle("/app/", middleware(http.StripPrefix("/app", http.FileServer(http.Dir("./web/")))))
}
