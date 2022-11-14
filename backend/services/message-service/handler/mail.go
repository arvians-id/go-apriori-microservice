package handler

import (
	"encoding/json"
	"github.com/arvians-id/go-apriori-microservice/services/message-service/config"
	"gopkg.in/gomail.v2"
	"log"
	"strconv"

	"github.com/nsqio/go-nsq"
)

type MailValue struct {
	ToEmail string
	Subject string
	Message string
}

type MailService struct {
	MailFromAdress string
	MailHost       string
	MailPort       string
	MailUsername   string
	MailPassword   string
}

func NewMailService(configuration *config.Config) *MailService {
	return &MailService{
		MailFromAdress: configuration.MailFromAddress,
		MailHost:       configuration.MailHost,
		MailPort:       configuration.MailPort,
		MailUsername:   configuration.MailUsername,
		MailPassword:   configuration.MailPassword,
	}
}

func (consumer *MailService) SendEmailWithText(message *nsq.Message) error {
	var mailValue MailValue
	if err := json.Unmarshal(message.Body, &mailValue); err != nil {
		log.Println("[EmailService][Unmarshal] unable to unmarshal data, err: ", err.Error())
		return err
	}

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", consumer.MailFromAdress)
	mailer.SetHeader("To", mailValue.ToEmail)
	mailer.SetHeader("Subject", mailValue.Subject)
	mailer.SetBody("text/html", mailValue.Message)

	port, err := strconv.Atoi(consumer.MailPort)
	if err != nil {
		log.Println("[EmailService][SendEmailWithText] problem in conversion string to integer, err: ", err.Error())
		return err
	}
	dialer := gomail.NewDialer(
		consumer.MailHost,
		port,
		consumer.MailUsername,
		consumer.MailPassword,
	)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		log.Println("[EmailService][SendEmailWithText] problem in mail sender, err: ", err.Error())
		return err
	}

	return nil
}
