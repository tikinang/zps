package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
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
		return
	}
	defer r.Body.Close()

	f, err := os.Create("books.db")
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

type book struct {
	authors string
	title   string
	notes   []string
}

func (h *handler) handleGet(w http.ResponseWriter, r *http.Request) {
	if !h.basicAuth(w, r) {
		return
	}

	db, err := sql.Open("sqlite3", "file:books.db")
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

	books := make(map[int]book)
	for rows.Next() {
		var oid int
		var authors, title, note string
		if err := rows.Scan(&oid, &authors, &title, &note); err != nil {
			log.Println("scan error:", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if b, has := books[oid]; has {
			b.notes = append(b.notes, note)
			books[oid] = b
		} else {
			books[oid] = book{
				authors: authors,
				title:   title,
				notes:   []string{note},
			}
		}
	}

	log.Println("print books.db notes")
	w.WriteHeader(http.StatusOK)
	for _, b := range books {
		fmt.Fprintf(w, "%s - %s\n\n", b.title, b.authors)
		for _, n := range b.notes {
			fmt.Fprintln(w, n)
		}
		fmt.Fprintln(w)
	}
}
