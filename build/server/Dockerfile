FROM golang:1.21 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o ./server ./cmd/server

FROM debian:stable

RUN apt-get update

RUN apt-get install -y libgomp1

RUN mkdir /app

ENV GO_ENV=production

COPY --from=builder /app/server /app

COPY --from=builder /app/pkg/keys /app/pkg/keys

WORKDIR /app

CMD ["./server"]