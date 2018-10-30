

FROM golang:latest

ARG telegram_bot_api_key
ENV TELEGRAM_BOT_API_KEY $telegram_bot_api_key

ARG listener_file
ENV LISTENER_FILE $listener_file

ENV GOPATH /go

WORKDIR /go/src/app

COPY ./main.go /go/src/app/main.go

ADD . .

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build -o main .

CMD ["./main"]