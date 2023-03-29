# Go as minimal as possible
FROM gcr.io/distroless/base-debian10

WORKDIR /

COPY ./calculator/calculator /calculator

EXPOSE 8080

USER root

ENTRYPOINT ["/calculator"]
