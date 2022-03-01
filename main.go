package main

import (
	"net/http"
	"std/internal/router"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Info("Starting local host server at :8000")
	http.ListenAndServe(":8000", router.RouterHandle())
}
