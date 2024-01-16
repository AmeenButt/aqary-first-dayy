FROM golang:latest

WORKDIR /app

COPY . .

RUN go get
RUN go build -o bin .

ENTRYPOINT [ "/app/bin" ]
