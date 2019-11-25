FROM golang:1.12

ENV GO111MODULE=on
ENV GOOS=linux GARCH=amd64 CGO_ENABLED=0

WORKDIR ${GOPATH}/src/github.com/adobe/prometheus-emcisilon-exporter

COPY . .

RUN go build

FROM alpine:3.10

COPY --from=0 /go/src/github.com/adobe/prometheus-emcisilon-exporter/prometheus-emcisilon-exporter /usr/bin/prometheus-emcisilon-exporter

EXPOSE 8080

ENTRYPOINT [ "prometheus-emcisilon-exporter" ]