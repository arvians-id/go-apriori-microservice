FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go build -o message-service

CMD ./message-service