# syntax=docker/dockerfile:1

FROM golang:1.19-alpine3.15

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . ./

RUN go build -o /stack-app-bin

CMD [ "/stack-app-bin" ]