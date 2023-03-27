#/bin/bash

CWD=$(pwd)
DATA_DIR=$CWD/data/

# Run k8s commands

# Everything else
kubectl delete -f auth/
kubectl delete -f coordinator/
