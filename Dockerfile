FROM golang:alpine as builder

WORKDIR /app

COPY app/go.mod app/go.sum ./

RUN go mod download

COPY app/ .

RUN go build -o main .

FROM alpine:latest

RUN apk --no-cache add ca-certificates mysql-client

COPY --from=builder /app/main .

CMD ["./main"]
