package main

import (
	"bytes"
	"encoding/json"
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

type SESMessage struct {
	Content          SESMessageContent
	FromEmailAddress string
	Destination      SESMessageDestination
}

type SESMessageContent struct {
	Simple SESMessageSimple
}

type SESMessageSimple struct {
	Body    SESMessageBody
	Subject SESMessageTextField
}

type SESMessageBody struct {
	Text SESMessageTextField
}

type SESMessageDestination struct {
	ToAddresses []string
}

type SESMessageTextField struct {
	Charset string
	Data    string
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

		form.Execute(w, data)

		if message.Address != "" && message.Subject != "" && message.Text != "" {
			sendMessage(message)
		}
	})

	http.ListenAndServe(":80", nil)
}

func sendMessage(message EmailMessage) {
	resp, err := handleAwsSes(message)
	log.Print(resp.StatusCode)
	if err != nil {
		log.Print(err)
	}

	/*response, err := handleSendGrid(message)
	if err != nil || response.StatusCode != 200 {

	}*/
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

func handleAwsSes(message EmailMessage) (*http.Response, error) {
	client := &http.Client{}

	to := message.Address
	subject := url.QueryEscape(message.Subject)
	text := url.QueryEscape(message.Text)
	from := os.Getenv("SENDER_ADDRESS")

	body := SESMessage{
		Content: SESMessageContent{
			Simple: SESMessageSimple{
				Body: SESMessageBody{
					Text: SESMessageTextField{
						Charset: "UTF-8",
						Data:    text,
					},
				},
				Subject: SESMessageTextField{
					Charset: "UTF-8",
					Data:    subject,
				},
			},
		},
		FromEmailAddress: from,
		Destination: SESMessageDestination{
			ToAddresses: []string{to},
		},
	}

	var query = "https://email.eu-north-1.amazonaws.com/v2/email/outbound-emails"
	apiKey := "AWS JBOUeiHcgl53NQ1VwT0QDiOM1zKOm83IrYPCzw2L" //os.Getenv("AWS_SES_SECRET_ACCESS_KEY")
	json, _ := json.Marshal(body)

	req, _ := http.NewRequest("POST", query, bytes.NewBuffer(json))
	req.Header.Add("Authorization", apiKey)
	req.Header.Add("Accept", "*/*")

	return client.Do(req)
}
