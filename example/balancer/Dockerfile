# Go as minimal as possible
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY ./balancer/balancer /balancer

EXPOSE 8080

USER root

ENTRYPOINT ["/balancer"]
