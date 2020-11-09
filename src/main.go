package main

import (
	"fmt"
	"jzio/src/internal/driver"
	"jzio/src/internal/handler"
	"net/http"
	"os"

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

	r.Get("/post/{id:[0-9]+}", handler.GetPost(c, cache))
	r.Get("/posts", handler.GetPosts(c, cache))

	fmt.Println("Serving on " + port)
	http.ListenAndServe(port, r)
	defer c.SQL.Close()
}
