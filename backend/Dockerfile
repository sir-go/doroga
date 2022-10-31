FROM golang:1.17-alpine3.16 as builder
RUN apk add --no-cache build-base imagemagick-dev imagemagick

WORKDIR /go/src/app
COPY go.mod .
COPY go.sum .
RUN go mod download
COPY ./cmd .

ENV CGO_ENABLED=1
ENV CGO_CFLAGS_ALLOW=-Xpreprocessor
ENV GOPROXY=https://proxy.golang.org

RUN go build -a -installsuffix cgo -o backend ./main

FROM alpine:3.16
RUN apk add --no-cache tzdata ca-certificates imagemagick-dev imagemagick
COPY --from=builder /go/src/app/backend /opt/backend
RUN mkdir -p /var/doroga_static/o/ /var/doroga_static/t/
WORKDIR /opt
EXPOSE 8081/tcp
ENTRYPOINT ["/opt/backend"]
