package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

var mapURL map[string]string

func main() {
	r := mux.NewRouter()

	srv := http.Server{
		Handler:      r,
		Addr:         "localhost:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

func ShortenURL(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	longURL := r.Form.Get("url")

	mapURL[longURL] = longURL
	fmt.Fprintln(w, http.StatusText(http.StatusOK))
}

func RedirectTo(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	longURL, ok := mapURL[vars["url"]]
	if !ok {
		http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}
