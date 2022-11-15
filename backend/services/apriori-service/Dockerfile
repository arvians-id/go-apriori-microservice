FROM golang:1.17-alpine

WORKDIR /app

COPY . .

RUN go build -o apriori-service

CMD ./apriori-service