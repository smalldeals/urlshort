package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httprate"
	"github.com/google/uuid"
)

/*
Create an http server
GET request will generate a key using uuid?
Each uuid key will map to URL.
Another request will get resolved url and redirect
user to the site via html
*/

// cache finds us the string from the query
var (
	cache map[string]string
	PORT  = os.Getenv("PORT")
)

func generateKey() string {
	return uuid.New().String()[0:8]
}

func serveGenerateLink(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	URI := r.URL.Query().Get("url")
	key := generateKey()

	cache[key] = URI
	fmt.Fprintf(w, "%v mapped to %v", URI, key)
}

func serveGetLink(w http.ResponseWriter, r *http.Request) {
	key := chi.URLParam(r, "key")
	URI, ok := cache[key]
	fmt.Println(URI, ok, key)
	if !ok {
		http.Error(w, "key does not exist", http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, URI, http.StatusSeeOther)
}

func serveHomePage(w http.ResponseWriter, r *http.Request) {
	io.WriteString(w, "Home Page Reached")
}

func main() {
	cache = make(map[string]string)
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(httprate.LimitByIP(50, 1*time.Minute))

	r.Get("/", serveHomePage)
	r.Get("/api/generatelink", serveGenerateLink)
	r.Get("/api/getlink/{key}", serveGetLink)
	http.ListenAndServe(":8080", r)
}
