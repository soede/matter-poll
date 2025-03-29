FROM golang:alpine AS builder

WORKDIR /build

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

RUN go build -o main cmd/bot/main.go

FROM alpine

WORKDIR /build

COPY --from=builder /build/main /build/main

EXPOSE 8080

CMD ["/build/main"]
