package main

import (
	"fmt"
	"io/ioutil"
	"log"
	// "net/smtp"
	"strings"
	"time"

	"github.com/Shobhitdimri01/Bookings/internal/models"
	mail "github.com/xhit/go-simple-mail/v2"
)

func ListenforMail() {
	go func() {
		for {
			msg := <-app.MailChan
			//SMTPEmail(msg)
			sendmsg(msg)

		}

	}()
}

//Sending Mail via Go-Simple-Mail
func sendmsg(m models.MailData) {
	password := "exrcnmmmausvfbol"
 	gmail := "smtp.gmail.com"
	server := mail.NewSMTPClient()
	server.Host = gmail
	server.Port = 587
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second

	server.Username = m.From
	server.Password = password
	server.Encryption = mail.EncryptionSTARTTLS

	// Since v2.3.0 you can specified authentication type:
	// - PLAIN (default)
	// - LOGIN
	// - CRAM-MD5
	// - None
	// server.Authentication = mail.AuthPlain

	// Variable to keep alive connection
	server.KeepAlive = false

	// Timeout for connect to SMTP Server
	server.ConnectTimeout = 10 * time.Second

	// Timeout for send the data and wait respond
	server.SendTimeout = 10 * time.Second

	// Set TLSConfig to provide custom TLS configuration. For example,
	// to skip TLS verification (useful for testing):
	// server.TLSConfig = &tls.Config{InsecureSkipVerify: true}


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



//Sending mail via Simple SMTP Server
//SMTP Server
// func SMTPEmail(m models.MailData) {
// 	password := "exrcnmmmausvfbol"
// 	gmail := "smtp.gmail.com"
// 	auth := smtp.PlainAuth(
// 		"",
// 		m.From,
// 		password,
// 		gmail,
// 	)

// 	headers := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";"
// 	var msgTosend string
// 	if m.Template == "" {
// 		msgTosend = m.Content
// 	} else {
// 		data, err := ioutil.ReadFile(fmt.Sprintf("./email-templates/%s", m.Template))
// 		if err != nil {
// 			app.ErrorLog.Println("Error in reading file", err)
// 		}

// 		mailTemplate := string(data)

// 		msgTosend = strings.Replace(mailTemplate, "[%body%]", m.Content, 1)
// 		// email.SetBody(mail.TextHTML, msgTosend)

// 		// html :="<h1>hello</h1>"
// 		msg := "Subject:" + m.Subject + "\n" + headers + "\n\n" + msgTosend

// 		GmailAddress:="smtp.gmail.com:587"

// 		err = smtp.SendMail(
// 			GmailAddress,
// 			auth,
// 			m.From,
// 			[]string{m.To},
// 			[]byte(msg),
// 		)

// 		if err != nil {
// 			fmt.Println(err)
// 		}
// 	}
// }
