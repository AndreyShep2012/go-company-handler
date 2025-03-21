FROM golang:1.23.2-alpine as builder

RUN apk add --no-cache make git

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN make build

FROM alpine:latest

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /app/bin/main .
COPY --from=builder /app/config.yml .

EXPOSE 8080

CMD ["./main"]
