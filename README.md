### GoMail

A simple mail service created for a Go code challenge. Designed to send e-mails using a primary service, with a secondary service as a fallback.

#### Components

The service consists of two components:
- a simple form page handling input and providing feedback
- a backend handling the underlying logic

Currently, the service uses [SendGrid](https://docs.sendgrid.com/v2-api) as the primary mailing provider and [Mailgun](https://documentation.mailgun.com/en/latest/index.html) as a fallback option.

### Running in Docker

In the `./src/gomail.env` file, configure the API keys and sender addresses for SendGrid and Mailgun, and the port the application is to be run on.

Technically, only the SendGrid configuration is required, but the fallback will not work without also configuring Mailgun.

To start the application, run the following commands in the terminal:

```
docker build --tag docker-go-mail .
docker run --publish [PORT]:[PORT] --env-file ./src/gomail.env docker-go-mail
```

### Running outside Docker

The application can be started from the command line, like any other Go application. 

Currently, there are no external dependencies used in the project. However, the application relies on environment variables to store API keys and sender addresses for the primary and fallback service, and these must be configured before running, along with the application port. See `./src/gomail.env` for a list of variables that need to be configured.

After configuring the environment, simply run `go run .` in the `src` folder.

### Testing

To test the application in case of provider failure, set an incorrect API key or sender address. Setting incorrect data for the primary will cause the application to use the fallback.
If both are misconfigured, an error message will be displayed in the UI. After each failure, the response status and body will be logged to the console.