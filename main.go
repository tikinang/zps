package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"
	"zps/pkg/api"
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

	h := api.NewHandler()
	r := mux.NewRouter()
	r.Path("/get/{key}").HandlerFunc(h.HandleGet)
	r.Path("/create/{key}/{value}").HandlerFunc(h.HandleCreate)
	r.Path("/remove/{key}").HandlerFunc(h.HandleRemove)
	r.Path("/list").HandlerFunc(h.HandleList)
	r.PathPrefix("/").HandlerFunc(h.HandleIndex)
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
