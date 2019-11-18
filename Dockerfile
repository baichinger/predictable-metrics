FROM golang:1.13.4 AS build-env
WORKDIR /go/src/predictable-metrics
COPY . .
RUN go build

FROM registry.access.redhat.com/ubi8/ubi-minimal:latest
COPY --from=build-env /go/src/predictable-metrics /app/
ENTRYPOINT /app/predictable-metrics
