package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/Shobhitdimri01/Bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func ListenforMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendmsg(msg)
		}

	}()
}

func sendmsg(m models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "localhost"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	// receiver_name := strings.Split(m.To, "@")
	if m.Template == "" {
		email.SetBody(mail.TextHTML, m.Content+" !")
	} else {
		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
		if err != nil {
			app.ErrorLog.Println("Error in reading file", err)
		}

		mailTemplate := string(data)

		msgTosend := strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
		email.SetBody(mail.TextHTML, msgTosend)
	}

	err = email.Send(client)
	if err != nil {
		log.Println(err.Error())
	} else {
		log.Println("Email sent !")
	}
}
