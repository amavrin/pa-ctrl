package k8s

import (
	"context"
	"sync"

	zlog "github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

const (
	CM_NAME = "hpa-config"
)

func WatchForChanges(ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
	mutex *sync.Mutex) {
	for {
		watcher, err := clientset.CoreV1().ConfigMaps(namespace).
			Watch(ctx,
				metav1.SingleObject(metav1.ObjectMeta{
					Name:      CM_NAME,
					Namespace: namespace,
				},
				),
			)
		if err != nil {
			panic(err)
		}
		updatePAConfig(watcher.ResultChan(), mutex)
	}
}

func updatePAConfig(eventChannel <-chan watch.Event, mutex *sync.Mutex) {
	for {
		event, open := <-eventChannel
		if open {
			switch event.Type {
			case watch.Added:
				fallthrough
			case watch.Modified:
				zlog.Print("watch.Modified")
				mutex.Lock()
				// Update our endpoint
				zlog.Print("TODO: update PA config here...")
				mutex.Unlock()
			case watch.Deleted:
				zlog.Print("watch.Deleted")
				mutex.Lock()
				// Fall back to the default value
				zlog.Print("TODO: update PA config here...")
				mutex.Unlock()
			default:
				zlog.Print("unknown event type")
				// Do nothing
			}
		} else {
			// If eventChannel is closed, it means the server has closed the connection
			return
		}
	}
}
