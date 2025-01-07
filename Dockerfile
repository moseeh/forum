FROM golang:1.22.2

LABEL maintainer="Moses Onyango ,Aaron Ochieng ,Swabri Musa ,Andy Osindo"
LABEL version="1.0"
LABEL description="An advanced golang web project"


COPY ./assets ./assets
COPY ./cmd ./cmd
COPY ./internal ./internal
COPY ./Makefile ./
COPY go.mod ./
COPY go.sum ./


RUN go mod tidy