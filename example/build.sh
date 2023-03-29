#/bin/bash

# Build
export GOOS=linux
export GOARCH=arm

go build -o ./balancer/balancer ./balancer/*.go
go build -o ./calculator/calculator ./calculator/*.go

kubectl delete -f example.yml

# Create images
eval $(minikube docker-env)
docker build --platform linux/arm64 -t ms-calculator:0.0.0 -f ./balancer/Dockerfile .
docker image push ms-calculator:0.0.0

docker build --platform linux/arm64 -t shard-calculator:0.0.0 -f ./calculator/Dockerfile .
docker image push shard-calculator:0.0.0

# Apply to cluster
kubectl apply -f example.yml
