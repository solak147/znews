package service

import (
	"net/smtp"
	"os"
	"znews/app/middleware"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

func Send(title string, body string, to string) {
	envErr := godotenv.Load()
	if envErr != nil {
		panic(envErr)
	}

	from := os.Getenv("GMAIL_USERNAME")
	pass := os.Getenv("GMAIL_PASSWORD")
	port := os.Getenv("GMAIL_PORT")
	server := os.Getenv("GMAIL_SERVER")

	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + title + "\n" +
		body

	err := smtp.SendMail(server+":"+port,
		smtp.PlainAuth("", from, pass, server),
		from, []string{to}, []byte(msg))
	if err != nil {
		middleware.Logger().WithFields(logrus.Fields{
			"title": "SendMail Failed",
		}).Error(err.Error())
		return
	}

}
