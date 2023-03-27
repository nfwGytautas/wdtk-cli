#/bin/bash

# Get pod id
PODS="$(kubectl get pod)"
AUTH_DB_POD=$(kubectl get pod | grep 'auth-db*' | awk '{print $1}')

if [ "$AUTH_DB_POD" = "" ]; then
    echo "Minikube not set up correctly or the pod is not running"
    exit
fi

# tty
kubectl exec --stdin --tty $AUTH_DB_POD -- /bin/bash
