FROM golang:1.20
WORKDIR /app

COPY src/go.mod ./
COPY src/*.go ./
COPY src/res/*.html ./res/

RUN CGO_ENABLED=0 GOOS=linux go build -o /docker-go-mail

CMD [ "/docker-go-mail" ]