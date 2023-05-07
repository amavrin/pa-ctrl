package config

import (
	"fmt"
	"io/ioutil"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

const (
	NS_FILE = "/var/run/secrets/kubernetes.io/serviceaccount/namespace"
	CM_NAME = "hpa-config"
)

// Load returns clientset and current namespace
func Load() (*kubernetes.Clientset, string, error) {
	nsBytes, err := ioutil.ReadFile(NS_FILE)
	if err != nil {
		return &kubernetes.Clientset{},
			"",
			fmt.Errorf("Unable to read namespace file at %s: %w", NS_FILE, err)
	}

	namespace := string(nsBytes)

	clientCfg, err := rest.InClusterConfig()
	if err != nil {
		return &kubernetes.Clientset{},
			"",
			fmt.Errorf("unable to get client configuration: %w", err)
	}

	clientset, err := kubernetes.NewForConfig(clientCfg)
	if err != nil {
		return clientset,
			"",
			fmt.Errorf("Unable to create our clientset: %w", err)
	}

	return clientset, namespace, nil
}
