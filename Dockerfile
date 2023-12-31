FROM golang:1.19-alpine as builder
WORKDIR /build

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .
RUN go build -o bin/web .

EXPOSE 4001

FROM alpine:latest
COPY --from=builder /build/bin .
COPY .env .
CMD ["./web"]