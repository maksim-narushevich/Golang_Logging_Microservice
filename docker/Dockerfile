FROM golang:latest

RUN apt-get update && apt-get install vim -y

RUN go get "github.com/joho/godotenv"
RUN go get "github.com/streadway/amqp"

ADD ./src /go/src/github.com/go-logging-micro-service

# Set working directory
WORKDIR /go/src/github.com/go-logging-micro-service
