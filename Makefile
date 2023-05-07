build:
	go build -o bin/pa-ctrl cmd/pa-ctrl/main.go
	go build -o bin/nursed-test cmd/nursed-test/main.go

image:
	docker build . -f Dockerfile.base -t pa-ctrl:base
	docker build . -t pa-ctrl:t1

image-nursed-test:
	docker build . -f Dockerfile.base -t pa-ctrl:base
	docker build . -f Dockerfile.nursed-test -t nursed-test:t1

load:
	kind load docker-image pa-ctrl:t1
	kind load docker-image nursed-test:t1

deploy:
	kubectl apply -f resources/ns.yaml
	kubectl -n test apply -f resources/pa-ctrl/rbac.yaml
	kubectl -n test apply -f resources/pa-ctrl/deploy.yaml

deploy-nursed-test:
	kubectl apply -f resources/ns.yaml
	kubectl -n test apply -f resources/nursed-test/deployment1.yaml
	kubectl -n test apply -f resources/nursed-test/deployment2.yaml

restart:
	kubectl -n test rollout restart deployment pa-ctrl

pf:
	kubectl -n test port-forward deployment/pa-ctrl 8080:8080

logs:
	stern -n test pa-ctrl

cm:
	kubectl -n test apply -f resources/pa-ctrl/cm.yaml

create-cluster:
	KUBECONFIG=~/.kube/config-colima kind create cluster
	kind get kubeconfig > ~/.kube/config-kind
	chmod 600 ~/.kube/config-kind

delete-cluster:
	KUBECONFIG=~/.kube/config-colima kind delete cluster

metrics-server:
	kubectl apply -k resources/metrics-server/
