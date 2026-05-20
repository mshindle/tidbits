FROM mshindle/golang-extra:latest AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/tidbits main.go

FROM debian:bookworm-slim

# install the necessary build packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates

RUN addgroup -S nonroot && adduser -S nonroot -G nonroot

WORKDIR /app
COPY --from=builder /app/tidbits /app/tidbits

USER nonroot:nonroot
CMD ["/app/tidbits"]
