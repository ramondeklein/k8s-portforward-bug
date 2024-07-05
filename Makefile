test: create-cluster docker deploy-k8s.yaml
	kind load docker-image test-server:latest
	-kubectl delete -f deploy-k8s.yaml
	sleep 3
	kubectl apply -f deploy-k8s.yaml
	sleep 3
	kubectl port-forward server 8080:8080

server: server.go
	GOOS=linux CGO_ENABLED=0 go build -o server server.go

docker: server Dockerfile
	docker build -t test-server:latest .

create-cluster:
	-kind create cluster
