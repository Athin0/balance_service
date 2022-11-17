FROM golang:latest

WORKDIR /app
COPY ./ /app

#install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

RUN chmod +x wait-for-postgres.sh

RUN go build -o myapp ./cmd/main.go

CMD ["/app/myapp"]