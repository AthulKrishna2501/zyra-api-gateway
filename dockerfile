FROM golang:1.24

WORKDIR /app

COPY . .

RUN go mod tidy

RUN go build -o api-gateway ./cmd

EXPOSE 5000

CMD ["./api-gateway"]
