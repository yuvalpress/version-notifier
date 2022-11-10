FROM golang:1.19.3-alpine

WORKDIR /app

COPY ./src ./

RUN go mod download
RUN go build -o /version-notifier

CMD [ "/version-notifier" ]