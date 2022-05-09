package main

import (
	"context"
	"database/sql"
	"embed"
	_ "embed"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"regexp"
	"sort"
	"strconv"
	"time"
	"zps/pkg/graceful"

	"github.com/gorilla/mux"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	ctx, cancel := graceful.Context()
	defer cancel()

	h := newHandler()
	r := mux.NewRouter()
	r.PathPrefix("/").Methods(http.MethodPut).HandlerFunc(h.handlePut)
	r.PathPrefix("/").Methods(http.MethodGet).HandlerFunc(h.handleGet)
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", 2922),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal("listen error:", err)
		}
	}()
	log.Println("server started")

	<-ctx.Done()
	log.Println("context done")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatal("server shutdown failed:", err)
	}
	log.Println("server shutdown")
}

type handler struct {
	password string
	username string
}

func newHandler() *handler {
	h := new(handler)
	var ok bool
	h.username, ok = os.LookupEnv("HTTP_BASIC_AUTH_USERNAME")
	if !ok {
		log.Fatal("HTTP_BASIC_AUTH_USERNAME env not found")
	}
	h.password, ok = os.LookupEnv("HTTP_BASIC_AUTH_PASSWORD")
	if !ok {
		log.Fatal("HTTP_BASIC_AUTH_PASSWORD env not found")
	}
	return h
}

func (h *handler) basicAuth(w http.ResponseWriter, r *http.Request) bool {
	username, password, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	if username != h.username || password != h.password {
		w.WriteHeader(http.StatusUnauthorized)
		return false
	}

	return true
}

func (h *handler) handlePut(w http.ResponseWriter, r *http.Request) {
	if !h.basicAuth(w, r) {
		log.Println("unauthorized request")
		return
	}
	defer r.Body.Close()

	f, err := os.Create("/mnt/storage/books.db")
	if err != nil {
		log.Println("create books.db file error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, r.Body)
	if err != nil {
		log.Println("copy to books.db file error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	log.Println("books.db updated")
	w.WriteHeader(http.StatusOK)
}

//go:embed query.sql
var query string

//go:embed list.html
var tpl embed.FS

type Book struct {
	Authors string
	Title   string
	Notes   []Note
}

type Note struct {
	Page int
	Text string
}

var pageRegex = regexp.MustCompile(`pbr:/word\?page=(\d*)`)

func (h *handler) handleGet(w http.ResponseWriter, r *http.Request) {

	t, err := template.ParseFS(tpl, "list.html")
	if err != nil {
		log.Println("parse template error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	db, err := sql.Open("sqlite3", "file:/mnt/storage/books.db")
	if err != nil {
		log.Println("open books.db error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer db.Close()

	rows, err := db.Query(query)
	if err != nil {
		log.Println("query error:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	books := make(map[string]Book)
	for rows.Next() {
		var authors, title, text, pageInfo string
		if err := rows.Scan(&authors, &title, &text); err != nil {
			log.Println("scan error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		var page int
		matches := pageRegex.FindStringSubmatch(pageInfo)
		if len(matches) == 2 {
			page, _ = strconv.Atoi(matches[1])
		}

		id := fmt.Sprintf("%s_%s", title, authors)
		if b, has := books[id]; has {
			note := Note{
				Page: page,
				Text: text,
			}
			b.Notes = append(b.Notes, note)
			books[id] = b
		} else {
			books[id] = Book{
				Authors: authors,
				Title:   title,
				Notes: []Note{{
					Page: page,
					Text: text,
				}},
			}
		}
	}

	log.Println("print books.db notes")
	w.WriteHeader(http.StatusOK)

	var sorted []Book
	for _, b := range books {
		sorted = append(sorted, b)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i].Title > sorted[j].Title
	})

	_ = t.Execute(w, map[string]interface{}{"Books": sorted})
}
