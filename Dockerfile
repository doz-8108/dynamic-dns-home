FROM golang:latest

WORKDIR /go-ddns

COPY . .

RUN go mod download
RUN go build -o app .

CMD ["./app"]