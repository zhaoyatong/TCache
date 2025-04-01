package main

import (
	"log"
	"net/http"
	"tcache/api"
)

func main() {
	addr := "127.0.0.1:9999"
	peers := api.NewHTTPPool(addr)
	log.Println("TCache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}
