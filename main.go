package main

import (
	"fmt"
	"god/proxy"
	"log"
)

func main() {
	fmt.Println("Hello, G.O.D!")

	// Proxy server start (temp)
	// Forward Proxy command (window) : curl -x http://127.0.0.1:8083 http://httpbin.org/get
	if err := proxy.Start("127.0.0.1:8083"); err != nil {
		log.Fatalf("Proxy faild: %v", err)
	}
}
