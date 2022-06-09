package api

import (
	"context"
	"crypto/sha512"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
)

type Handler struct {
	rdb *redis.Client
}

func NewHandler() *Handler {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("db_hostname"), os.Getenv("db_port")),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Printf("connected to redis: %s\n", cn)
			return nil
		},
		DB: 1,
	})

	return &Handler{
		rdb: rdb,
	}

}

func (h *Handler) Close() {}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 1024*1024)
	unixNano := time.Now().UnixNano()
	random := rand.New(rand.NewSource(unixNano))
	if _, err := random.Read(b); err != nil {
		fmt.Fprintf(w, "error reading rand: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	hash := sha512.Sum512(b)
	key := strconv.FormatInt(unixNano, 10)
	stringHash := fmt.Sprintf("%x", hash)
	setCmd := h.rdb.Set(r.Context(), key, stringHash, 0)
	if err := setCmd.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	getCmd := h.rdb.Get(r.Context(), key)
	if err := getCmd.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "%v\n", err)
		return
	}
	if getCmd.Val() != stringHash {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "redis value differs: %s != %s\n", getCmd.Val(), stringHash)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "handled with timestamp: %s\n", time.Now())
	fmt.Fprintf(w, "%x\n", hash)
}
