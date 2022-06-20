# syntax=docker/dockerfile:1

FROM golang:1.18.3-alpine as builder
WORKDIR /usr/src/app
COPY go.mod go.sum ./
RUN go mod download && go mod verify
COPY . .
RUN go build -v -o /usr/local/bin/evbus-server ./cmd/server/

FROM alpine:3.16.0
COPY --from=builder /usr/local/bin/evbus-server /usr/local/bin/evbus-server
EXPOSE 3366
ENTRYPOINT ["evbus-server"]
