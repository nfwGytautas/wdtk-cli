# Backends

This directory contains the backend bindings for creating microservices in here you will find directories for
languages that provide various bindings for creating microservices in different languages

## Supported languages

### Golang (go/)

Support for both balancer and microservice api

## NOTE
Altho currently only golang is supported the design of mstk doesn't prohibit the use of any language for a balancer backend
because the only thing you need for it is a simple http server that can receive and handle http requests. While the services can be made in
any different supported backend communication protocol
