package main

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

const defaultPort = ":8080"

const static = "/swagger"
const docdir = "/swagger/doc/"

func main() {
	r := mux.NewRouter()
	r.PathPrefix(static).Handler(http.FileServer(http.Dir(".")))
	r.PathPrefix("/swagger/*").Handler(fileserver())

	n := negroni.New()
	n.Use(recovery())
	n.UseHandler(r)

	n.Run(defaultPort)
}

// rec .
type rec struct{}

// recovery .
func recovery() *rec {
	return new(rec)
}

// ServeHTTP .
func (rec) ServeHTTP(rw http.ResponseWriter, r *http.Request, next http.HandlerFunc) {
	defer func() {
		if err := recover(); err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
			data, _ := json.Marshal(err)
			rw.Write(data)
		}
	}()

	next(rw, r)
}

type fs string

// fileserver .
func fileserver() *fs {
	return new(fs)
}

// ServeHTTP .
func (fs fs) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(r.URL.Path))
}
