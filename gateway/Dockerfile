FROM golang:1.21 as builder

WORKDIR /app
COPY go.mod go.mod
COPY go.sum go.sum
COPY go.work go.work
COPY go.work.sum go.work.sum
COPY api api
COPY gateway gateway
WORKDIR gateway
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /usr/local/bin/gateway

FROM gcr.io/distroless/static-debian11
COPY --from=builder /usr/local/bin/gateway /usr/local/bin/gateway
ENTRYPOINT ["/usr/local/bin/gateway"]
