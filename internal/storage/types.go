package storage

import "sync"

type Target struct {
	Name         string `yaml:"name"`
	UpperCPULoad int    `yaml:"upperCPULoad"`
	Min          int    `yaml:"minReplicas,omitempty"`
	Max          int    `yaml:"maxReplicas,omitempty"`
}

type TargetStatus struct {
	Name            string `yaml:"name"`
	Found           bool   `yaml:"found"`
	CurrentReplicas int    `yaml:"currentReplicas"`
	DeriredReplicas int    `yaml:"desiredReplicas"`
	CurrentCPULoad  int    `yaml:"currentCPULoad"`
}

var Mu sync.Mutex

type Data struct {
	Targets []Target
	Status  []TargetStatus
}
