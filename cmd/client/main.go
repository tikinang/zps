package main

import (
	"log"
	"net/http"
	"os"
)

var (
	Username      string
	Password      string
	ServerAddress string
)

func main() {
	f, err := os.Open("/mnt/ext1/system/config/books.db")
	if err != nil {
		log.Fatalf("could not open books.db file: %v\n", err)
	}
	defer f.Close()

	r, err := http.NewRequest(http.MethodPut, ServerAddress, f)
	if err != nil {
		log.Fatalf("failed to create http request: %v\n", err)
	}
	r.SetBasicAuth(Username, Password)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatalf("http request error: %v\n", err)
	}
	defer resp.Body.Close()

	log.Println("response status code:", resp.StatusCode)
}
