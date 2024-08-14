FROM golang:1.22 AS build_base

WORKDIR /tmp/dogstatsd-local

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 go test -v ./...

RUN CGO_ENABLED=0 GOOS=linux go build -o ./out/dogstatsd-local ./cmd/dogstatsd-local/main.go

FROM scratch

COPY --from=build_base /tmp/dogstatsd-local/out/dogstatsd-local /app/dogstatsd-local

EXPOSE 8125

ENTRYPOINT ["/app/dogstatsd-local"]
CMD ["/app/dogstatsd-local"]