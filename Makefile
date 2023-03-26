build:
	go build -o pa-ctrl cmd/main.go

image:
	docker build . -t pa-ctrl:t1

load:
	kind load docker-image pa-ctrl:t1

deploy:
	kubectl apply -f resources/ns.yaml
	kubectl -n test apply -f resources/rbac.yaml
	kubectl -n test apply -f resources/deploy.yaml

restart:
	kubectl -n test rollout restart deployment pa-ctrl

pf:
	kubectl -n test port-forward deployment/pa-ctrl 8080:8080

logs:
	stern -n test pa-ctrl

cm:
	kubectl -n test apply -f resources/cm.yaml

cluster:
	kind create cluster