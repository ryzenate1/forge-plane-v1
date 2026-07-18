FROM golang:1.26-alpine AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/daemon ./cmd/daemon

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
RUN mkdir -p /srv/game-panel/servers
RUN adduser -D -H -u 10002 appuser
USER appuser
COPY --from=build /out/daemon /daemon
EXPOSE 9090 2022
HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 CMD ["/daemon", "--healthcheck"]
ENTRYPOINT ["/daemon"]
