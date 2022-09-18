# syntax=docker/dockerfile:1

FROM golang:1.19-alpine3.15 as build

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /stack-app-bin

FROM alpine:latest

WORKDIR /app
COPY --from=build /stack-app-bin /stack-app-bin

ENTRYPOINT [ "/stack-app-bin" ]