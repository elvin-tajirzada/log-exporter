FROM golang:1.20-alpine

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go mod verify
RUN GOOS=linux go build -o ./bin/log-exporter ./cmd/log-exporter

ENTRYPOINT /app/bin/log-exporter