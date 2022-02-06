package api

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

type Handler struct {
	random *rand.Rand
	rdb    *redis.Client
}

func NewHandler() *Handler {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("redis0_hostname"), os.Getenv("redis0_port")),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Printf("connected to redis: %s", cn)
			return nil
		},
		DB: 1,
	})

	return &Handler{
		random: rand.New(rand.NewSource(time.Now().Unix())),
		rdb:    rdb,
	}
}

func (h *Handler) HandleIndex(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 4096)
	if _, err := h.random.Read(b); err != nil {
		fmt.Println("error reading rand:", err)
	}
	hash := sha512.Sum512(b)

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "handled with timestamp: %s\n\n", time.Now())
	fmt.Fprintf(w, "%x\n", hash)
}

func (h *Handler) HandleGet(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	cmd := h.rdb.Get(r.Context(), key)
	if err := cmd.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "get: %s -> %s\n", key, cmd.Val())
}

func (h *Handler) HandleCreate(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key, val := vars["key"], vars["value"]

	cmd := h.rdb.Set(r.Context(), key, val, 0)
	if err := cmd.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "created: %s -> %s\n", key, val)
}

func (h *Handler) HandleRemove(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	key := vars["key"]

	cmd := h.rdb.Del(r.Context(), key)
	if err := cmd.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "removed: %s\n", key)
}

func (h *Handler) HandleList(w http.ResponseWriter, r *http.Request) {
	cmd := h.rdb.Keys(r.Context(), "*")
	if err := cmd.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%+v\n", err)
		return
	}

	w.WriteHeader(http.StatusOK)
	for _, key := range cmd.Val() {

		var val string
		cmd := h.rdb.Get(r.Context(), key)
		if err := cmd.Err(); err != nil {
			val = err.Error()
		} else {
			val = cmd.Val()
		}

		fmt.Fprintf(w, "list: %s -> %s\n", key, val)
	}
}
