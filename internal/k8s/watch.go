package k8s

import (
	"context"

	"github.com/amavrin/pa-ctrl/internal/storage"
	zlog "github.com/rs/zerolog/log"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/kubernetes"
)

const (
	CM_NAME    = "hpa-config"
	CONFIG_KEY = "deployments.yaml"
)

func WatchForChanges(ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string) {
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
		err = updatePAConfig(watcher.ResultChan())
		if err != nil {
			zlog.Print("cannot update config: ", err)
		}
	}
}

func updatePAConfig(eventChannel <-chan watch.Event) error {
	for {
		event, open := <-eventChannel
		if open {
			switch event.Type {
			case watch.Added:
				fallthrough
			case watch.Modified:
				zlog.Print("watch.Modified")
				storage.Mu.Lock()
				// Update our endpoint
				zlog.Print("updating PA config...")
				if updatedMap, ok := event.Object.(*corev1.ConfigMap); ok {
					if config, ok := updatedMap.Data[CONFIG_KEY]; ok {
						cfg, err := storage.SaveTargets(config)
						if err != nil {
							return err
						}
						storage.Config = *cfg
					}
				}
				storage.Mu.Unlock()
			case watch.Deleted:
				zlog.Print("watch.Deleted")
				storage.Mu.Lock()
				// Fall back to the default value
				zlog.Print("TODO: update PA config here...")
				storage.Mu.Unlock()
			default:
				zlog.Print("unknown event type")
				// Do nothing
			}
		} else {
			// If eventChannel is closed, it means the server has closed the connection
			return nil
		}
	}
}
