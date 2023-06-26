# GoMail

A simple mail service created for a Go code challenge. Designed to send e-mails using a primary service, with a secondary service as a fallback.

## Approach

The application is built to be as simple as possible, as I was also effectively learning Go as I went along. The input is handled by a simple form in a Go HTML template.
The application is hardcoded to use [SendGrid](https://docs.sendgrid.com/v2-api) as the primary mailing provider and [Mailgun](https://documentation.mailgun.com/en/latest/index.html) as a fallback option - if the first fails, it will log an error message along with the HTTP status code and response body. In case of the fallback also failing, another log will be made and a message will be displayed in the UI.

The application requires API keys and sender addresses to be provided either through environment variables or an .env file passed to the Docker container.

As investigating automated testing mostly pointed to third-party libraries which I was trying to avoid, the application has to be tested manually. This can be done by passing incorrect config data, to either trigger the fallback or a complete failure.

## Running in Docker

In the `./src/gomail.env` file, configure the API keys and sender addresses for SendGrid and Mailgun, and the port the application is to be run on.

Technically, only the SendGrid configuration is required, but the fallback will not work without also configuring Mailgun.

To start the application, run the following commands in the terminal:

```
docker build --tag docker-go-mail .
docker run --publish [PORT]:[PORT] --env-file ./src/gomail.env docker-go-mail
```

## Running outside Docker

The application can be started from the command line, like any other Go application. 

Currently, there are no external dependencies used in the project. However, the application relies on environment variables to store API keys and sender addresses for the primary and fallback service, and these must be configured before running, along with the application port. See `./src/gomail.env` for a list of variables that need to be configured.

After configuring the environment, simply run `go run .` in the `src` folder.