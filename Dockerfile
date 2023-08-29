FROM mshindle/golang-extra:latest AS builder

# copy our application code over...
COPY . /app
WORKDIR /app
RUN go mod download
RUN go build -o /tmp/tidbits main.go

FROM debian:bullseye-slim

# install the necessary build packages
RUN apt-get update && \
    apt-get install -y --no-install-recommends ca-certificates
WORKDIR /app
COPY --from=builder /tmp/tidbits /app/tidbits
ENTRYPOINT ["/app/tidbits"]
