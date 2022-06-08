package api

import (
	"context"
	"crypto/sha512"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/go-redis/redis/v8"
)

type Handler struct {
	rdb   *redis.Client
	files []io.WriteCloser
	index *uint64
}

func NewHandler() *Handler {
	var files []io.WriteCloser
	for i := 0; i < 16; i++ {
		f, err := os.Create(fmt.Sprintf("stressor_file_%d", time.Now().UnixNano()))
		if err != nil {
			log.Fatalf("failed to create file: %v\n", err)
		}
		files = append(files, f)
	}

	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("db_hostname"), os.Getenv("db_port")),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Printf("connected to redis: %s\n", cn)
			return nil
		},
		DB: 1,
	})

	return &Handler{
		rdb:   rdb,
		files: files,
		index: new(uint64),
	}

}

func (h *Handler) Close() {
	for _, f := range h.files {
		if err := f.Close(); err != nil {
			log.Printf("error closing file: %v\n", err)
		}
	}
}

func (h *Handler) Handle(w http.ResponseWriter, r *http.Request) {
	b := make([]byte, 1024*1024)
	unixNano := time.Now().UnixNano()
	random := rand.New(rand.NewSource(unixNano))
	if _, err := random.Read(b); err != nil {
		fmt.Fprintf(w, "error reading rand: %v\n", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	i := atomic.AddUint64(h.index, 1) % 16
	if _, err := h.files[i].Write(b); err != nil {
		fmt.Fprintf(w, "error writing to file: %v\n", err)
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
