package main

import "proxy/http"

func main() {
	http.Listen(":8080")
}