package main

import (
	"crypto/subtle"
	"net/http"
	"os"

	"github.com/explabs/prometheus-manager/routers"
)

func BasicAuth(handler http.HandlerFunc, username, password, realm string) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		user, pass, ok := r.BasicAuth()

		if !ok || subtle.ConstantTimeCompare([]byte(user), []byte(username)) != 1 || subtle.ConstantTimeCompare([]byte(pass), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="`+realm+`"`)
			w.WriteHeader(401)
			w.Write([]byte("Unauthorised.\n"))
			return
		}

		handler(w, r)
	}
}

func main() {
	username := "admin"
	password := os.Getenv("ADMIN_PASS")
	http.HandleFunc("/start", BasicAuth(routers.StartContainer, username, password, ""))
	http.HandleFunc("/stop", BasicAuth(routers.StopContainert, username, password, ""))
<<<<<<< HEAD
	http.HandleFunc("/generate", BasicAuth(routers.JsonParse, username, password, ""))
=======
	http.HandleFunc("/generate", BasicAuth(routers.GenerateConfigHandler, username, password, ""))
>>>>>>> 0b6ef4013ac27b6fe2a2d5e351b7bdc481df3907
	http.ListenAndServe(":9091", nil)
}
