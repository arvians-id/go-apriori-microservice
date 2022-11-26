package main

import (
	"github.com/arvians-id/go-apriori-microservice/services/message-service/config"
	"github.com/arvians-id/go-apriori-microservice/services/message-service/handler"
	"github.com/arvians-id/go-apriori-microservice/services/message-service/messaging"
	"log"
	"os"
	"os/signal"
	"syscall"
)

const (
	channelMail = "mail_channel"
	topicMail   = "mail_topic"
)

func main() {
	configuration, err := config.LoadConfig()
	if err != nil {
		log.Fatalln("Failed at config", err)
	}

	mailHandler := handler.NewMailService(configuration)
	mailConsumer := messaging.ConsumerConfig{
		Topic:         topicMail,
		Channel:       channelMail,
		LookupAddress: "nsq-release-nsqlookupd-0.nsq-release-nsqlookupd:4161",
		MaxAttempts:   10,
		MaxInFlight:   100,
		Handler:       mailHandler.SendEmailWithText,
	}

	mail := messaging.NewConsumer(mailConsumer)
	mail.Run()

	// keep app alive until terminated
	term := make(chan os.Signal, 1)
	signal.Notify(term, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	select {
	case <-term:
		log.Println("Application terminated")
	}
}
