# Kubes

This directory contains kubernetes clusters for development of MSTK

The development is done inside a local ```minikube``` cluster.

## Setup
For UNIX type systems a ```setup.sh``` is provided that automatically sets the environment up for you

For Windows users command ```kubectl apply -f {{dirs}}``` on all directories inside ```kubes```

NOTE: That files with ```*-storage.yml``` have a special token {{pv_mount_point}} which should be replaced with the directory where you want
to store your minikube PVs, the setup script automatically sets this to the current directory
