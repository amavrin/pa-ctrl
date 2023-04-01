package k8s

import (
	"context"
	"reflect"
	"testing"

	appsv1 "k8s.io/api/apps/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"

	"k8s.io/client-go/kubernetes/fake"
)

func TestGetTargetDeployment(t *testing.T) {
	type args struct {
		ctx            context.Context
		clientset      kubernetes.Interface
		namespace      string
		deploymentName string
	}
	tests := []struct {
		name    string
		args    args
		want    appsv1.Deployment
		wantErr bool
	}{
		{
			name: "Single deployment in the correct namespace",
			args: args{
				ctx: context.Background(),
				clientset: fake.NewSimpleClientset(&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "deployment1",
						Namespace: "namespace1",
					},
				}),
				namespace:      "namespace1",
				deploymentName: "deployment1",
			},
			want: appsv1.Deployment{
				ObjectMeta: metav1.ObjectMeta{
					Name:      "deployment1",
					Namespace: "namespace1",
				},
			},
		},
		{
			name: "No deployment in the correct namespace",
			args: args{
				ctx: context.Background(),
				clientset: fake.NewSimpleClientset(&appsv1.Deployment{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "deployment1",
						Namespace: "namespace2",
					},
				}),
				namespace:      "namespace1",
				deploymentName: "deployment1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetTargetDeployment(tt.args.ctx,
				tt.args.clientset,
				tt.args.namespace,
				tt.args.deploymentName)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSubjectDeployment() error = %v, wantErr %v",
					err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSubjectDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}
