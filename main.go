package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"time"
	"zps/pkg/api"
	"zps/pkg/graceful"
)

func main() {
	ctx, cancel := graceful.Context()
	defer cancel()

	h := api.NewHandler()
	defer h.Close()
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(h.Handle)
	srv := &http.Server{
		Handler:      r,
		Addr:         fmt.Sprintf(":%d", 8080),
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
