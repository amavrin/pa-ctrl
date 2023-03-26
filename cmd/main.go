package main

import (
	"context"
	"os"
	"sync"

	"github.com/alexey.mavrin/pa-ctrl/cmd/internal/config"
	"github.com/alexey.mavrin/pa-ctrl/internal/k8s"
	"github.com/alexey.mavrin/pa-ctrl/internal/server"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	mutex := &sync.Mutex{}

	clientset, namespace, err := config.Load()
	if err != nil {
		zlog.Print("configure error: ", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go k8s.WatchForChanges(ctx, clientset, namespace, mutex)

	go server.Serve(mutex)

	podName := os.Getenv("POD_NAME")
	lockName := "pa-ctrl-lock"
	zlog.Print("getting lock ", lockName, ", identity ", podName,
		", namespace ", namespace)
	lock := k8s.GetNewLock(clientset, lockName, podName, namespace)
	k8s.RunLeaderElection(ctx, lock, podName)
}
