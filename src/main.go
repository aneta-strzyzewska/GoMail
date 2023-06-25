package main

import (
	"encoding/base64"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
)

type PageData struct {
	PageTitle string
	Message   EmailMessage
	Status    string
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
			Address: strings.TrimSpace(r.FormValue("email")),
			Subject: strings.TrimSpace(r.FormValue("subject")),
			Text:    strings.TrimSpace(r.FormValue("text")),
		}

		var status string
		if message.Address != "" && message.Subject != "" && message.Text != "" {
			status = sendMessage(message)
		}

		data := PageData{
			PageTitle: "GoMail",
			Message:   message,
			Status:    status,
		}

		form.Execute(w, data)
	})

	http.ListenAndServe(":"+os.Getenv("GOMAIL_PORT"), nil)
}

func sendMessage(message EmailMessage) string {
	resp, err := handleSendGrid(message)
	if err == nil && resp.StatusCode == 200 {
		return fmt.Sprintf("Message sent to %s", message.Address)
	}

	logError(resp)
	log.Println("Attempting to send via fallback")

	resp, err = handleMailgun(message)
	if err == nil && resp.StatusCode == 200 {
		return fmt.Sprintf("Message sent via fallback service to %s", message.Address)
	}

	logError(resp)
	log.Println("Failed to send message via fallback")

	return "Failed to send message"
}

func handleSendGrid(message EmailMessage) (*http.Response, error) {
	client := &http.Client{}

	to := message.Address
	subject := url.QueryEscape(message.Subject)
	text := url.QueryEscape(message.Text)
	from := url.QueryEscape(fmt.Sprintf("Go Mail <%s>", os.Getenv("PRIMARY_SENDER_ADDRESS")))

	query := fmt.Sprintf("https://api.sendgrid.com/api/mail.send.json?to=%s&subject=%s&text=%s&from=%s", to, subject, text, from)
	apiKey := "Bearer " + os.Getenv("SENDGRID_API_KEY")

	req, _ := http.NewRequest("POST", query, nil)
	req.Header.Add("Authorization", apiKey)
	req.Header.Add("Accept", "*/*")

	return client.Do(req)
}

func handleMailgun(message EmailMessage) (*http.Response, error) {
	client := &http.Client{}

	sender := os.Getenv("FALLBACK_SENDER_ADDRESS")
	to := message.Address
	subject := url.QueryEscape(message.Subject)
	text := url.QueryEscape(message.Text)
	from := url.QueryEscape(fmt.Sprintf("Go Mail <mailgun@%s>", sender))

	query := fmt.Sprintf("https://api.mailgun.net/v3/%s/messages?to=%s&subject=%s&text=%s&from=%s", sender, to, subject, text, from)
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

func logError(resp *http.Response) {
	log.Println("Could not send email")
	log.Println(resp.Status)
	body, _ := io.ReadAll(resp.Body)
	log.Println(string(body))
}
