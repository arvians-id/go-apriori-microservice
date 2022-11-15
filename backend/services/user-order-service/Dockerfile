FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go build -o user-order-service

CMD ./user-order-service