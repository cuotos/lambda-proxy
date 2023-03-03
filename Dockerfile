FROM golang:alpine as builder

WORKDIR /app

COPY go.mod go.sum .

RUN go mod download

COPY . .

RUN go build -o /lambda-proxy

FROM alpine

COPY --from=builder /lambda-proxy /lambda-proxy

ENTRYPOINT ["/lambda-proxy"]
