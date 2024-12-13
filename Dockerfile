# syntax=docker/dockerfile:1

FROM golang:1.23.2 as builder

WORKDIR /todo

COPY go.mod go.sum ./

COPY vendor/ ./vendor

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o todo ./cmd/todo/

FROM alpine:latest

LABEL maintainer="moazedy@gmail.com"

WORKDIR /root/

COPY --from=builder /todo/todo .
COPY --from=builder /todo/config/ ./config/ 

EXPOSE 4853

CMD ./todo
