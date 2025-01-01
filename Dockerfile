# Basis-Image für Go
FROM golang:1.20-alpine

WORKDIR /app

COPY go.mod ./
RUN go mod tidy

COPY . .

EXPOSE 8080

RUN go build -o main .

CMD ["./main"]
