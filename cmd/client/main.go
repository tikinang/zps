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
		log.Fatal("could not open books.db file:", err)
	}
	defer f.Close()

	r, err := http.NewRequest(http.MethodPut, ServerAddress, f)
	if err != nil {
		log.Fatal("failed to create http request", err)
	}
	r.SetBasicAuth(Username, Password)
	resp, err := http.DefaultClient.Do(r)
	if err != nil {
		log.Fatal("http request error:", err)
	}
	defer resp.Body.Close()

	log.Println("response status code:", resp.StatusCode)
}
