package main

import (
	"html/template"
	"net/http"
)

type PageData struct {
	PageTitle string
	Message   EmailMessage
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

		data := PageData{
			PageTitle: "GoMail",
			Message:   message,
		}

		form.Execute(w, data)
	})

	http.ListenAndServe(":80", nil)
}
