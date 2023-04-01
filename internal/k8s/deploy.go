package k8s

import (
	"context"

	appsv1 "k8s.io/api/apps/v1"
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
