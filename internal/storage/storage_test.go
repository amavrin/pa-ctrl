package storage

import (
	"reflect"
	"testing"
)

func TestSaveTargets(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name    string
		args    args
		want    *Data
		wantErr bool
	}{
		{
			name: "single deployment config",
			args: args{
				s: "- name: dep1",
			},
			want: &Data{
				Targets: []Target{
					{
						Name: "dep1",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := SaveTargets(tt.args.s)
			if (err != nil) != tt.wantErr {
				t.Errorf("SaveTargets() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SaveTargets() = %v, want %v", got, tt.want)
			}
		})
	}
}
