FROM golang:latest

WORKDIR /app
COPY go.* ./
RUN go mod download

COPY ./ /app

RUN apt-get update && apt-get -y upgrade

RUN go build -o myapp ./cmd/main.go

CMD ["/app/myapp"]