package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

var (
	// httpPort port
	httpPort = ":8888"
	// swaggerAPI api folder
	swaggerAPI = "api"
)

func init() {
	if envPort := os.Getenv("SWAGGER_PORT"); len(envPort) != 0 {
		httpPort = envPort
	}
	if envAPI := os.Getenv("SWAGGER_DOC"); len(envAPI) != 0 {
		swaggerAPI = envAPI
	}

	os.MkdirAll(swaggerAPI, 0755)
}

func main() {
	r := mux.NewRouter()
	r.PathPrefix("/swagger").Handler(fileserver())

	n := negroni.New()
	n.Use(recovery())
	n.UseHandler(r)

	n.Run(httpPort)
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

type fs struct {
	fileServer    http.Handler
	apiFileServer http.FileSystem
}

// fileserver .
func fileserver() *fs {
	return &fs{
		fileServer:    http.FileServer(http.Dir(".")),
		apiFileServer: http.Dir(swaggerAPI),
	}
}

// ServeHTTP .
func (fs *fs) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, ".json") {
		api := strings.TrimPrefix(r.URL.Path, "/swagger/")

		f, err := fs.apiFileServer.Open(api)
		if err != nil {
			http.Error(rw, fmt.Sprintf("%s not found or permission denied", api), http.StatusNotFound)
			return
		}

		var buf = new(bytes.Buffer)
		io.Copy(buf, f)
		f.Close()

		rw.Header().Set("Content-Type", "text/html; charset=utf-8")
		rw.WriteHeader(http.StatusOK)
		rw.Write(buf.Bytes())
		return
	}

	fs.fileServer.ServeHTTP(rw, r)
}
