FROM golang:latest
LABEL authors="nikitakeba"

WORKDIR /app

COPY configs /app/configs
COPY . .

RUN go build -o driver-service ./cmd/app/main.go

CMD ["./driver-service"]
