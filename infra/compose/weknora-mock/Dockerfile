# Mock WeKnora service for CI/integration testing
# Provides minimal endpoints needed for Servify integration

FROM alpine:latest

RUN apk --no-cache add ca-certificates

# Install a simple HTTP server (netcat-openbsd or use busybox httpd)
# For simplicity, we'll use a basic Go binary compiled in a temp location

COPY --from=golang:1.23-alpine /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:$PATH

WORKDIR /app

# Create a simple Go server
RUN cat > main.go << 'EOF'
package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

func main() {
	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"service": "weknora-mock",
			"version": "1.0.0",
		})
	})

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status":  "ok",
			"message": "mock weknora service",
		})
	})

	fmt.Println("Mock WeKnora service listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", nil))
}
EOF

RUN go mod init weknora-mock && \
    go get -d && \
    CGO_ENABLED=0 go build -o weknora-mock .

EXPOSE 9000 9001

CMD ["./weknora-mock"]
