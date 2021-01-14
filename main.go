package main

import (
	"log"
	"net/http"
)

func main() {
	app := NewApplication()
	log.Fatal(http.ListenAndServe(":8080", app.Router))
}
