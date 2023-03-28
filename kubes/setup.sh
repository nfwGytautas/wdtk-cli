#/bin/bash

CWD=$(pwd)
DATA_DIR=$CWD/data

# Make sure we are inside kubes directory
SANITY_CHECK=$(basename $CWD)
if [ "$SANITY_CHECK" != "kubes" ]; then
    echo "Run setup.sh inside kubes directory"
    exit
fi

# Run k8s commands
echo "Setting up minikube kubernetes cluster"

# Everything else
kubectl apply -f auth/
kubectl apply -f coordinator/
kubectl apply -f gateway/
