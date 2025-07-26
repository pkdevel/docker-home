FROM --platform=$BUILDPLATFORM golang:1.24-alpine AS fetcher
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=bind,source=go.mod,target=go.mod \
  --mount=type=bind,source=go.sum,target=go.sum \
  go mod download -x

FROM ghcr.io/a-h/templ:latest AS templ
WORKDIR /app
COPY --chown=65532:65532 web web
RUN ["templ", "generate"]

FROM node:alpine AS tailwindcss
WORKDIR /app
RUN --mount=type=bind,source=package.json,target=package.json \
  ["npm", "install"]
RUN --mount=type=bind,source=web,target=web \
  ["npx", "@tailwindcss/cli", "-i", "web/template/tailwind.css", "-o", "assets/style.css", "-m"]

FROM fetcher AS builder
WORKDIR /app
ARG TARGETARCH
ENV CGO_ENABLED=0
RUN --mount=type=cache,target=/go/pkg/mod \
  --mount=type=bind,target=. \
  --mount=type=bind,source=app/web,target=web,from=templ \
  GOARCH=$TARGETARCH go build -v -o /bin/docker-home ./cmd/main.go

FROM scratch AS assembler
EXPOSE 8080
LABEL name="docker-home"
LABEL description="A simple docker home page"
WORKDIR /
COPY assets assets
COPY --from=tailwindcss app/assets/. assets/
COPY --from=builder /bin/docker-home bin/
ENTRYPOINT ["docker-home"]
