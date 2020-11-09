package handler

import (
	"context"
	"database/sql"
	"strconv"
	//	"errors"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	lru "github.com/hashicorp/golang-lru"
	"jzio/src/internal/driver"
)

type Post struct {
	id              int64
	Url             string `json:"url"`
	Title           string `json:"title"`
	Thumbnail       string `json:thumbnail"`
	Likes           int32  `json:likes"`
	Views           int32  `json:views"`
	CreatedAt       int32  `json:timestamp"`
	DurationSeconds int32  `json:duration_seconds"`
}

// GetPosts() GET returns a User Vault
// request with a user_public_id
func GetPosts(db *driver.DB, cache *lru.ARCCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		tx, err := db.SQL.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		// start and end Unixtime stamps
		posts, err := loadPosts(ctx, tx, cache)

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		respondwithJSON(w, http.StatusCreated, posts)
	}
}

// GetPost() GET returns a User Vault
// request with a user_public_id
func GetPost(db *driver.DB, cache *lru.ARCCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		ctx := context.Background()

		tx, err := db.SQL.BeginTx(ctx, nil)
		if err != nil {
			log.Fatal(err)
		}

		id, err := strconv.Atoi(chi.URLParam(r, "id"))
		if err != nil {
			log.Fatal(err)
		}

		// start and end Unixtime stamps
		posts, err := loadPost(ctx, tx, cache, id)

		err = tx.Commit()
		if err != nil {
			log.Fatal(err)
		}

		respondwithJSON(w, http.StatusCreated, posts)
	}
}

// loadPost() Loads a single vault record by user_id
func loadPost(ctx context.Context, tx *sql.Tx, cache *lru.ARCCache, id int) (*Post, error) {
	query := "SELECT id, url, title, thumbnail, likes, views, created_at, duration_seconds FROM posts WHERE id=? LIMIT 1"

	p := Post{}
	if value, ok := cache.Get(id); ok {
		log.Print("hit cache")
		return value.(*Post), nil
	} else {
		log.Print("missed cache")
		rows, err := tx.QueryContext(
			ctx,
			query,
			id,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			switch err := rows.Scan(&p.id, &p.Url, &p.Title, &p.Thumbnail, &p.Likes, &p.Views, &p.CreatedAt, &p.DurationSeconds); err {
			case sql.ErrNoRows:
			case nil:
				cache.Add(id, &p)
				return &p, nil
			default:
				panic(err)
			}
		}
	}

	return &p, nil
}

// loadPosts() Loads a single vault record by user_id
func loadPosts(ctx context.Context, tx *sql.Tx, cache *lru.ARCCache) ([]*Post, error) {
	query := "SELECT id, url, title, thumbnail, likes, views, created_at, duration_seconds FROM posts LIMIT 5"
	posts := make([]*Post, 0)

	keys := cache.Keys()
	if len(keys) > 0 {
		log.Print("hit cache")
		for _, v := range keys {
			value, ok := cache.Get(v.(int))
			if ok {
				posts = append(posts, value.(*Post))
			}
		}
	} else {

		rows, err := tx.QueryContext(
			ctx,
			query,
		)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			p := Post{}
			switch err := rows.Scan(&p.id, &p.Url, &p.Title, &p.Thumbnail, &p.Likes, &p.Views, &p.CreatedAt, &p.DurationSeconds); err {
			case sql.ErrNoRows:
			case nil:
				posts = append(posts, &p)
			default:
				panic(err)
			}
		}
	}
	return posts, nil
}
