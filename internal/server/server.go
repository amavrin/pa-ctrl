package server

import (
	"fmt"
	"net/http"
	"sync"

	zlog "github.com/rs/zerolog/log"
)

const (
	configStub = "some configuration data"
)

func Serve(mutex *sync.Mutex) {
	http.HandleFunc("/config", func(w http.ResponseWriter, r *http.Request) {
		mutex.Lock()
		body := []byte(fmt.Sprintf(`{"current_config": "%s"}`, configStub))
		mutex.Unlock()
		w.WriteHeader(http.StatusOK)
		w.Write(body)
	})

	zlog.Print("Listening on port 8080")
	http.ListenAndServe(":8080", nil)
}
