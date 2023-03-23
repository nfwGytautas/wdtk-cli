# gomods
In this directory you'll find the various go modules that can be used when creating microservices.

## Stand alone services

### ```gateway```
Module containing functionality needed to create an API gateway (Includes authentication functions).

### ```coordinator```
The coordinator of MSTK (This is an app that needs to run on a server).

### ```balancers```
Various load balancers provided by MSTK

## APIs/Helpers

### ```coordinator-api```
The common API for MSTK coordinator. Other packages link to this one to unify communication with a coordinator.

### ```balancer-api```
The common API for MSTK balancers.

### ```common```
Common functions used by packages

### ```microservice-api```
Functions used for creating a microservice, end users will use this library the most
