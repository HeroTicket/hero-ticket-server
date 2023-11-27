FROM golang:1.21-alpine3.18 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o ./subscriber ./cmd/subscriber

FROM alpine:3.18

ENV GO_ENV=production

RUN mkdir /app

COPY --from=builder /app/swagger /app

WORKDIR /app

CMD ["./subscriber"]
