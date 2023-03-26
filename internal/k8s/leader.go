package k8s

import (
	"context"
	"time"

	zlog "github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

func GetNewLock(clientset *kubernetes.Clientset,
	lockname, podname, namespace string) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockname,
			Namespace: namespace,
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: podname,
		},
	}
}

func doStuff() {
	for {
		zlog.Print("doing stuff...")
		time.Sleep(5 * time.Second)
	}
}

func RunLeaderElection(ctx context.Context, lock *resourcelock.LeaseLock, id string) {
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
				doStuff()
			},
			OnStoppedLeading: func() {
				zlog.Print("no longer the leader, staying inactive.")
			},
			OnNewLeader: func(current_id string) {
				if current_id == id {
					zlog.Print("still the leader!")
					return
				}
				zlog.Print("new leader is ", current_id)
			},
		},
	})
}
