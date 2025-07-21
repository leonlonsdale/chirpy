package main

import "net/http"

func (app *application) createFileServer() http.Handler {
	return http.StripPrefix("/app", http.FileServer(http.Dir("./web/")))
}
