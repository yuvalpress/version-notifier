FROM golang:1.19.3-alpine

WORKDIR /app

COPY ./src ./

ENV GOTRACEBACK "none"

RUN go mod download
RUN go build -o /version-notifier

RUN apk add chromium

CMD [ "/version-notifier" ]