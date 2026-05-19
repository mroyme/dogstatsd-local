FROM cgr.dev/chainguard/go:latest AS build_base

WORKDIR /tmp/dogstatsd-local

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o ./out/dogstatsd-local ./cmd/dogstatsd-local/main.go

FROM scratch

LABEL org.opencontainers.image.source="https://github.com/mroyme/dogstatsd-local" \
      org.opencontainers.image.description="A local DogStatsD protocol inspector for debugging metrics, service checks, and events" \
      org.opencontainers.image.licenses="MIT" \
      org.opencontainers.image.title="dogstatsd-local" \
      org.opencontainers.image.url="https://mroyme.github.io/dogstatsd-local/" \
      org.opencontainers.image.documentation="https://mroyme.github.io/dogstatsd-local/" \
      org.opencontainers.image.authors="Madhurjya Roy <m@mroy.me>" \
      org.opencontainers.image.vendor="Madhurjya Roy"

COPY --from=build_base /tmp/dogstatsd-local/out/dogstatsd-local /app/dogstatsd-local

EXPOSE 8125

ENTRYPOINT ["/app/dogstatsd-local"]
CMD ["/app/dogstatsd-local"]
