package k8s

import (
	"context"
	"sync"
	"time"

	"github.com/amavrin/pa-ctrl/internal/storage"
	zlog "github.com/rs/zerolog/log"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/leaderelection"
	"k8s.io/client-go/tools/leaderelection/resourcelock"
)

var Leader sync.Mutex

func GetNewLock(clientset *kubernetes.Clientset,
	lockname, id, namespace string) *resourcelock.LeaseLock {
	return &resourcelock.LeaseLock{
		LeaseMeta: metav1.ObjectMeta{
			Name:      lockname,
			Namespace: namespace,
		},
		Client: clientset.CoordinationV1(),
		LockConfig: resourcelock.ResourceLockConfig{
			Identity: id,
		},
	}
}

func ProcessTargetDeployments(ctx context.Context,
	clientset *kubernetes.Clientset,
	namespace string,
) {
	for {
		Leader.Lock()
		// get performance data and update replica count
		zlog.Print("updating deployment data...")

		storage.Mu.Lock()
		storage.Config.Status = make([]storage.TargetStatus, 0)
		for _, dep := range storage.Config.Targets {
			deploy, err := GetTargetDeployment(ctx,
				clientset,
				namespace,
				dep.Name,
			)
			if err != nil {
				zlog.Print(err)
				continue
			}
			storage.Config.Status = append(storage.Config.Status,
				storage.TargetStatus{Name: deploy.ObjectMeta.Name})
		}
		storage.Mu.Unlock()

		time.Sleep(5 * time.Second)
		Leader.Unlock()
	}
}

func RunLeaderElection(ctx context.Context,
	clientset *kubernetes.Clientset,
	id string,
	namespace string,
) {
	lockName := "pa-ctrl-lock"
	zlog.Print("getting lock ", lockName, ", identity ", id,
		", namespace ", namespace)
	lock := GetNewLock(clientset, lockName, id, namespace)
	leaderelection.RunOrDie(ctx, leaderelection.LeaderElectionConfig{
		Lock:            lock,
		ReleaseOnCancel: true,
		LeaseDuration:   15 * time.Second,
		RenewDeadline:   10 * time.Second,
		RetryPeriod:     2 * time.Second,
		Callbacks: leaderelection.LeaderCallbacks{
			OnStartedLeading: func(c context.Context) {
				zlog.Print("Leader, unlocking leader lock")
				Leader.Unlock()
			},
			OnStoppedLeading: func() {
				zlog.Print("no longer the leader, locking leader lock")
				Leader.Lock()
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
