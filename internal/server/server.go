package server

import (
	"fmt"
	"net/http"

	"github.com/amavrin/pa-ctrl/internal/storage"
	zlog "github.com/rs/zerolog/log"
	"gopkg.in/yaml.v2"
)

func printConfig(w http.ResponseWriter, r *http.Request) {
	storage.Mu.Lock()
	config, err := yaml.Marshal(storage.Config)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := []byte(fmt.Sprintf(`%s`, config))
	storage.Mu.Unlock()
	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

func Serve() {
	mux := http.NewServeMux()

	mux.HandleFunc("/config", printConfig)

	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	mux.HandleFunc("/ready", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	zlog.Print("Listening on port 8080")
	http.ListenAndServe(":8080", mux)
}
