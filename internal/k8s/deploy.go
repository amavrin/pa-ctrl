package k8s

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// GetTargetDeployment returns the deployment
// matched name specified
func GetTargetDeployment(ctx context.Context,
	clientset kubernetes.Interface,
	namespace string,
	deploymentName string) (appsv1.Deployment, error) {

	deployment, err := clientset.AppsV1().Deployments(namespace).Get(ctx, deploymentName, metav1.GetOptions{})
	if err != nil {
		return appsv1.Deployment{}, err
	}
	return *deployment, nil
}

func GetPodsForDeployment(ctx context.Context,
	clientset kubernetes.Interface,
	namespace string,
	deployment appsv1.Deployment,
) ([]corev1.Pod, error) {
	replicaSet, err := getReplicaSetForDeployment(ctx, clientset, namespace, deployment)
	if err != nil {
		return nil, err
	}

	pods, err := getPodsForReplicaSet(ctx, clientset, namespace, replicaSet)

	if err != nil {
		return nil, err
	}

	return pods, nil
}

func getReplicaSetForDeployment(ctx context.Context,
	clientset kubernetes.Interface,
	namespace string,
	deployment appsv1.Deployment,
) (appsv1.ReplicaSet, error) {
	replicaSets, err := clientset.AppsV1().ReplicaSets(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return appsv1.ReplicaSet{}, err
	}
	for _, rs := range replicaSets.Items {
		if rs.ObjectMeta.OwnerReferences[0].Kind != "Deployment" {
			continue
		}
		if rs.ObjectMeta.OwnerReferences[0].Name != deployment.ObjectMeta.Name {
			continue
		}
		return rs, nil
	}
	return appsv1.ReplicaSet{}, fmt.Errorf("No ReplicaSet found for deployment %s",
		deployment.ObjectMeta.Name)
}

func getPodsForReplicaSet(ctx context.Context,
	clientset kubernetes.Interface,
	namespace string,
	replicaSet appsv1.ReplicaSet,
) ([]corev1.Pod, error) {
	pods := make([]corev1.Pod, 0)
	return pods, nil
}
