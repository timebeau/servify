FROM golang:1.25-alpine AS builder

WORKDIR /app
# Build uses module in apps/server
WORKDIR /app/apps/server
COPY apps/server/go.mod .
COPY apps/server/go.sum .
RUN go mod download

WORKDIR /app
COPY . .
# Build server binary by default (apps/server)
RUN CGO_ENABLED=0 GOOS=linux go build -o servify ./apps/server/cmd/server

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

COPY --from=builder /app/servify .
# Config file can be mounted as ./config.yml at runtime

EXPOSE 8080

CMD ["./servify"]
