package router

import (
	"github.com/gorilla/mux"
	"io"
	"net/http"
)

func RegisterHandlers(r *mux.Router) {
	r.HandleFunc("/health", heatlhcheckHandler)
	r.HandleFunc("/", hello)
}

func heatlhcheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	io.WriteString(w, `{"alive": true}`)
}

func hello(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Hello, World!\n"))
}
