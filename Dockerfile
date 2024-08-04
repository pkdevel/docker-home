FROM golang:1.22-alpine3.20 AS fetcher
COPY go.mod go.sum /app/
WORKDIR /app
RUN go mod download

FROM ghcr.io/a-h/templ:latest AS templ
COPY --chown=65532:65532 ./web/. /app/web
WORKDIR /app
RUN ["templ", "generate"]

FROM d3fk/tailwindcss:latest AS tailwindcss
COPY ./assets/. /app/assets
COPY ./web/. /app/web
WORKDIR /app
RUN ["/tailwindcss", "-c", "web/tailwind.config.js", "-i", "web/template/tailwind.css", "-o", "assets/style.css", "-m"]

FROM fetcher AS builder
COPY ./. /app
COPY --from=templ /app/web/. /app/web
WORKDIR /app
RUN go build -v -o main ./cmd/main.go

FROM alpine:3.20 AS assembler
LABEL name="docker-home"
LABEL description="A simple docker home page"
EXPOSE 8080
COPY --from=tailwindcss /app/assets /assets
COPY --from=builder /app/main /docker-home
WORKDIR /
ENTRYPOINT ["/docker-home"]

