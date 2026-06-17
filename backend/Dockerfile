# Multi-stage, non-root, без :latest (tech-stack §7). Собирает server и seed.
FROM golang:1.23-alpine AS build
WORKDIR /src
RUN apk add --no-cache git
COPY . .
# go.sum не коммитим в этот MVP — генерируем при сборке (tech-stack: vendored снимает токен,
# здесь публичные модули, поэтому tidy в build достаточно).
RUN go mod tidy
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/server ./cmd/server
RUN CGO_ENABLED=0 go build -trimpath -ldflags "-s -w" -o /out/seed ./cmd/seed

FROM alpine:3.20
RUN apk add --no-cache ca-certificates wget && \
    adduser -D -u 10001 app && \
    mkdir -p /data/blobs && chown -R app /data
COPY --from=build /out/server /usr/local/bin/server
COPY --from=build /out/seed /usr/local/bin/seed
USER app
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/server"]
