package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"log"
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
	Text    string
}

func main() {
	form := template.Must(template.ParseFiles("res/form.html"))
	http.HandleFunc("/go-mail", func(w http.ResponseWriter, r *http.Request) {
		message := EmailMessage{
			Address: r.FormValue("email"),
			Subject: r.FormValue("subject"),
			Text:    r.FormValue("text"),
		}

		data := PageData{
			PageTitle: "GoMail",
			Message:   message,
		}

		if message.Address != "" && message.Subject != "" && message.Text != "" {
			sendMessage(message)
		}

		form.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}

func sendMessage(message EmailMessage) {
	resp, err := handleSendGrid(message)
	if err != nil || resp.StatusCode != 200 {
		log.Printf("Could not send email")
		log.Printf("Attempting to send via fallback")
		handleMailgun(message)
	}
}

func handleSendGrid(message EmailMessage) (*http.Response, error) {
	client := &http.Client{}

	to := message.Address
	subject := url.QueryEscape(message.Subject)
	text := url.QueryEscape(message.Text)
	from := os.Getenv("SENDER_ADDRESS")

	query := fmt.Sprintf("https://api.sendgrid.com/api/mail.send.json?to=%s&subject=%s&text=%s&from=%s", to, subject, text, from)
	apiKey := "Bearer " + os.Getenv("SENDGRID_API_KEY")

	req, _ := http.NewRequest("POST", query, nil)
	req.Header.Add("Authorization", apiKey)
	req.Header.Add("Accept", "*/*")

	return client.Do(req)
}

func handleMailgun(message EmailMessage) (*http.Response, error) {
	client := &http.Client{}

	to := message.Address
	subject := url.QueryEscape(message.Subject)
	text := url.QueryEscape(message.Text)
	from := url.QueryEscape(fmt.Sprintf("Go Mail <%s>", "mailgun@sandbox1f13a792e65646d7ac6126cbcf158f5f.mailgun.org"))

	query := fmt.Sprintf("https://api.mailgun.net/v3/sandbox1f13a792e65646d7ac6126cbcf158f5f.mailgun.org/messages?to=%s&subject=%s&text=%s&from=%s", to, subject, text, from)
	apiKey := os.Getenv("MAILGUN_API_KEY")

	req, _ := http.NewRequest("POST", query, nil)
	req.Header.Add("Authorization", "Basic "+basicAuth("api", apiKey))
	req.Header.Add("Accept", "*/*")

	return client.Do(req)
}

func basicAuth(username, password string) string {
	auth := username + ":" + password
	return base64.StdEncoding.EncodeToString([]byte(auth))
}
