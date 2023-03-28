#/bin/bash

CWD=$(pwd)
DATA_DIR=$CWD/data/

# Run k8s commands, ignore the result of the error cause it should only be an error of 'not found'
kubectl delete -f auth/ || true
kubectl delete -f coordinator/ || true
kubectl delete -f gateway/ || true
