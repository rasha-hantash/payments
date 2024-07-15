FROM golang:1.21 as builder

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
COPY go.work go.work
COPY go.work.sum go.work.sum
COPY api api
COPY gateway gateway
WORKDIR api
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/api

FROM gcr.io/distroless/static-debian11
COPY --from=builder /usr/local/bin/api /usr/local/bin/api
ENTRYPOINT ["/usr/local/bin/api"]
