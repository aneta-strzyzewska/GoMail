### GoMail

A simple mail service created for a code challenge. Designed to send e-mails using a primary service, with a secondary service as a fallback.

#### Components

The service consists of three components:
- a simple frontend, built using Go templating
- a backend handling most of the logic

Currently, the service uses SendGrid as the primary mailing provider and Mailgun as a fallback option.

### Running in Docker

In the `./src/gomail.env` file, configure the API keys and sender addresses for SendGrid and Mailgun, and the port the application is to be run on.

Technically, only the SendGrid config is required, but the fallback will not work without also configuring Mailgun.

To start the application, run the following commands in the terminal:

```
docker build --tag docker-go-mail .
docker run --publish [PORT]:[PORT] --env-file ./src/gomail.env docker-go-mail
```

### Running outside Docker

The application can be started from the command line, like any other Go application. 

Currently, there are no extrenal dependencies used in the project. However, the application relies on environment variables to store API keys and sender addresses for the primary and fallback service, and these must be configured before running, along with the application port. See `./src/gomail.env` for a list of variables that need to be configured.

After configuring the environment, simply run `go run .` in the `src` folder.
