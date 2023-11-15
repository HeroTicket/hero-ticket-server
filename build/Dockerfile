FROM golang:1.21 AS builder

RUN mkdir /app

WORKDIR /app

COPY . .

RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -a -o ./server ./cmd

FROM debian:stable

RUN apt-get update

RUN apt-get install -y libgomp1

RUN mkdir /app

COPY --from=builder /app/server /app

COPY --from=builder /app/keys /app/keys

WORKDIR /app

CMD ["./server"]
