package main

import (
	"fmt"
	"html/template"
	"net/http"
	"net/url"
	"os"
)

type PageData struct {
	PageTitle string
	Message   EmailMessage
	Error     string
}

type EmailMessage struct {
	Address string
	Subject string
	Message string
}

func main() {
	form := template.Must(template.ParseFiles("res/form.html"))
	http.HandleFunc("/go-mail", func(w http.ResponseWriter, r *http.Request) {
		message := EmailMessage{
			Address: r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Message: r.FormValue("message"),
		}

		client := &http.Client{}

		to := message.Address
		subject := url.QueryEscape(message.Subject)
		text := url.QueryEscape(message.Message)
		from := os.Getenv("SENDER_ADDRESS")

		query := fmt.Sprintf("https://api.sendgrid.com/api/mail.send.json?to=%s&subject=%s&text=%s&from=%s", to, subject, text, from)
		apiKey := "Bearer " + os.Getenv("SENDGRID_API_KEY")

		req, _ := http.NewRequest("POST", query, nil)
		req.Header.Add("Authorization", apiKey)
		req.Header.Add("Accept", "*/*")

		data := PageData{
			PageTitle: "GoMail",
			Message:   message,
		}

		client.Do(req)

		form.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}
