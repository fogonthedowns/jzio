package main

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"jzio/src/internal/driver"
	"jzio/src/internal/handler"

	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
)

func main() {
	port := os.Getenv("SERVER_PORT")
	dbName := os.Getenv("DB_NAME")
	dbPass := os.Getenv("DB_PASS")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")

	cache, _ := lru.NewARC(3)

	r := chi.NewRouter()
	var err error
	c, err := driver.ConnectSQL(dbHost, dbPort, "root", dbPass, dbName)
	if err != nil {
		os.Exit(-1)
	}

	workDir, _ := os.Getwd()
	filesDir := http.Dir(filepath.Join(workDir, "data"))
	assetsDir := http.Dir(filepath.Join(workDir, "assets"))
	imgDir := http.Dir(filepath.Join(workDir, "img"))

	FileServer(r, "/assets", assetsDir)
	FileServer(r, "/img", imgDir)
	FileServer(r, "/", filesDir)

	r.Get("/post/{id:[0-9]+}", handler.GetPost(c, cache))
	r.Get("/posts.json", handler.GetPosts(c, cache))

	fmt.Println("Serving on " + port)
	http.ListenAndServe(port, r)
	defer c.SQL.Close()
}

// FileServer conveniently sets up a http.FileServer handler to serve
// static files from a http.FileSystem.
func FileServer(r chi.Router, path string, root http.FileSystem) {
	if strings.ContainsAny(path, "{}*") {
		panic("FileServer does not permit any URL parameters.")
	}

	if path != "/" && path[len(path)-1] != '/' {
		r.Get(path, http.RedirectHandler(path+"/", 301).ServeHTTP)
		path += "/"
	}
	path += "*"
	r.Get(path, func(w http.ResponseWriter, r *http.Request) {
		rctx := chi.RouteContext(r.Context())
		pathPrefix := strings.TrimSuffix(rctx.RoutePattern(), "/*")
		fmt.Println(pathPrefix)
		fs := http.StripPrefix(pathPrefix, http.FileServer(root))
		fs.ServeHTTP(w, r)
	})
}
