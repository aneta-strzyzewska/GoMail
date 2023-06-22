package main

import (
	"fmt"
	"html/template"
	"io"
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

		to := "gretelostrich@gmail.com"
		subject := url.QueryEscape("Test email")
		text := url.QueryEscape("testing the email")
		from := "gretelostrich@gmail.com"

		query := fmt.Sprintf("https://api.sendgrid.com/api/mail.send.json?to=%s&subject=%s&text=%s&from=%s", to, subject, text, from)
		log.Println(query)

		apiKey := "Bearer " + os.Getenv("SENDGRID_API_KEY")
		log.Print("apiKey", apiKey)

		req, _ := http.NewRequest("POST", query, nil)
		req.Header.Add("Authorization", apiKey)
		req.Header.Add("Accept", "*/*")

		data := PageData{
			PageTitle: "GoMail",
			Message:   message,
		}

		resp, _ := client.Do(req)
		log.Println("status: ", resp.Status)
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatal(err)
		}
		bodyString := string(bodyBytes)
		log.Println(bodyString)

		form.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}
