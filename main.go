package main

import (
	"net/http"

	"github.com/explabs/prometheus-manager/routers"
)

func main() {
	http.HandleFunc("/start", routers.StartContainer)
	http.HandleFunc("/stop", routers.StopContainert)

	http.ListenAndServe(":9091", nil)
}
