#/bin/bash

# Execute a bin push on minikube docker

eval $(minikube docker-env)
cd ../../bin
python3 build.py push
