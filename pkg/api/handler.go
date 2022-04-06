package api

import (
	"context"
	"crypto/sha512"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

type Handler struct {
	random *rand.Rand
	rdb    *redis.Client
}

func NewHandler() *Handler {
	rdb := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("db_hostname"), os.Getenv("db_port")),
		OnConnect: func(ctx context.Context, cn *redis.Conn) error {
			log.Printf("connected to db: %s", cn)
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
	env := os.Environ()
	var jsonEnv map[string]string
	if raw, has := os.LookupEnv("JSON_ENV"); has {
		_ = json.Unmarshal([]byte(raw), &jsonEnv)
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "handled with timestamp: %s\n", time.Now())
	fmt.Fprintf(w, "%x\n\n", hash)
	fmt.Fprintf(w, "os.Environ():\n")
	for _, e := range env {
		fmt.Fprintf(w, "%s\n", e)
	}
	if len(jsonEnv) > 0 {
		fmt.Fprintf(w, "\nJSON_ENV\n")
		for k, v := range jsonEnv {
			fmt.Fprintf(w, "%s=%s\n", k, v)
		}
	}
}
