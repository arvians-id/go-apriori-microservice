FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go build -o password-reset-service

CMD ./password-reset-service