package main

import (
	"context"
	"crypto/sha512"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
	"zps/pkg/graceful"
)

const ZpsPort = "ZPS_LISTEN_PORT"

func main() {
	ctx, cancel := graceful.Context()
	defer cancel()

	listenPort := 8080
	if envPort, has := os.LookupEnv(ZpsPort); has {
		listenPort, _ = strconv.Atoi(envPort)
	}

	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(handler)
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", listenPort),
		WriteTimeout: 5 * time.Second,
		ReadTimeout:  5 * time.Second,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen error: %s\n", err)
		}
	}()
	log.Println("server started")

	<-ctx.Done()
	log.Println("context done")

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown failed: %+v\n", err)
	}
	log.Println("server shutdown")
}

func handler(w http.ResponseWriter, r *http.Request) {

	random := rand.New(rand.NewSource(time.Now().Unix()))
	b := make([]byte, 512)
	if _, err := random.Read(b); err != nil {
		fmt.Println("error reading rand:", err)
	}
	hash := sha512.Sum512(b)

	w.Header().Set("Content-Type", "text/plain")
	w.WriteHeader(http.StatusOK)

	fmt.Fprintf(w, "handled with timestamp: %s\n\n", time.Now())
	fmt.Fprintf(w, "%x", hash)
}
