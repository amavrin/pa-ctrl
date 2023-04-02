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

var (
	deployment1 = appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment1",
			Namespace: "namespace1",
		},
	}

	replicaSet1 = appsv1.ReplicaSet{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "deployment1-xxxxxx",
			Namespace: "namespace1",
			OwnerReferences: []metav1.OwnerReference{
				{
					Kind: "Deployment",
					Name: "deployment1",
				},
			},
		},
	}
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
				ctx:            context.Background(),
				clientset:      fake.NewSimpleClientset(&deployment1),
				namespace:      "namespace1",
				deploymentName: "deployment1",
			},
			want: deployment1,
		},
		{
			name: "No deployment in the correct namespace",
			args: args{
				ctx:            context.Background(),
				clientset:      fake.NewSimpleClientset(),
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

func Test_getReplicaSetForDeployment(t *testing.T) {
	type args struct {
		ctx        context.Context
		clientset  kubernetes.Interface
		namespace  string
		deployment appsv1.Deployment
	}
	tests := []struct {
		name    string
		args    args
		want    appsv1.ReplicaSet
		wantErr bool
	}{
		{
			name: "Deployment has ReplicaSer",
			args: args{
				ctx:        context.Background(),
				clientset:  fake.NewSimpleClientset(&deployment1, &replicaSet1),
				namespace:  "namespace1",
				deployment: deployment1,
			},
			want: replicaSet1,
		},
		{
			name: "No ReplicaSet for Deployment in this namespace",
			args: args{
				ctx:       context.Background(),
				clientset: fake.NewSimpleClientset(&deployment1),
				namespace: "namespace1",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getReplicaSetForDeployment(tt.args.ctx, tt.args.clientset, tt.args.namespace, tt.args.deployment)
			if (err != nil) != tt.wantErr {
				t.Errorf("getReplicaSetForDeployment() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getReplicaSetForDeployment() = %v, want %v", got, tt.want)
			}
		})
	}
}
