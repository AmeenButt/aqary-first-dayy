FROM golang:latest

RUN mkdir /app
WORKDIR /app

COPY . .

RUN go run main.go

EXPOSE 8080

CMD ["./app"]
