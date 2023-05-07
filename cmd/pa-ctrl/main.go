package main

import (
	"context"
	"os"

	"github.com/amavrin/pa-ctrl/cmd/internal/config"
	"github.com/amavrin/pa-ctrl/internal/k8s"
	"github.com/amavrin/pa-ctrl/internal/server"
	zlog "github.com/rs/zerolog/log"
)

func main() {
	clientset, namespace, err := config.Load()
	if err != nil {
		zlog.Print("configure error: ", err)
		os.Exit(1)
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go k8s.WatchForChanges(ctx, clientset, namespace)

	go server.Serve()

	k8s.Leader.Lock()
	go k8s.ProcessTargetDeployments(ctx, clientset, namespace)

	podName := os.Getenv("POD_NAME")
	k8s.RunLeaderElection(ctx, clientset, podName, namespace)
}
