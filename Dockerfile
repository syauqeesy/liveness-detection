FROM golang:1.27.1-alpine3.24 AS builder

WORKDIR /app

COPY . .

RUN go mod download

RUN go build -o main .

FROM alpine:3.24 AS container

WORKDIR /app

COPY --from=builder /app/main .

EXPOSE 80

ENTRYPOINT ["./main"]

CMD ["http"]
