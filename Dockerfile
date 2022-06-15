# syntax=docker/dockerfile:1

FROM golang:1.18.3-alpine as build
WORKDIR /evbus
COPY . .
RUN go mod tidy && go mod vendor
RUN CGO_ENABLED=0 go build -a -o ./bin/ -tags netgo -ldflags '-w -extldflags "-static"' ./cmd/server/

FROM scratch
COPY --from=build /evbus/bin/server /evbus/bin/server
EXPOSE 3366
ENTRYPOINT ["/evbus/bin/server"]
