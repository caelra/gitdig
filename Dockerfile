# syntax = docker/dockerfile:1

FROM golang:alpine

WORKDIR /app

COPY . /app

RUN go build -o gitdig

ENTRYPOINT ["./gitdig"]
