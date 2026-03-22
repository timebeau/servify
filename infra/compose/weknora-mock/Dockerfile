# Mock WeKnora service for CI/integration testing
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=golang:1.23-alpine /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:$PATH

WORKDIR /app

# Create Go source file using printf to avoid heredoc issues
RUN printf '%s\n' \
'package main' \
'' \
'import (' \
'	"fmt"' \
'	"log"' \
'	"net/http"' \
'	"encoding/json"' \
')' \
'' \
'func main() {' \
'	http.HandleFunc("/api/v1/health", func(w http.ResponseWriter, r *http.Request) {' \
'		w.Header().Set("Content-Type", "application/json")' \
'		json.NewEncoder(w).Encode(map[string]interface{}{' \
'			"status":  "ok",' \
'			"service": "weknora-mock",' \
'			"version": "1.0.0",' \
'		})' \
'	})' \
'	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {' \
'		w.Header().Set("Content-Type", "application/json")' \
'		json.NewEncoder(w).Encode(map[string]interface{}{' \
'			"status":  "ok",' \
'			"message": "mock weknora service",' \
'		})' \
'	})' \
'	fmt.Println("Mock WeKnora service listening on :9000")' \
'	log.Fatal(http.ListenAndServe(":9000", nil))' \
'}' > main.go

RUN go mod init weknora-mock && \
    go get -d && \
    CGO_ENABLED=0 go build -o weknora-mock .

EXPOSE 9000 9001

CMD ["./weknora-mock"]
