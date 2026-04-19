# Mock Dify service for CI/integration testing
FROM alpine:latest

RUN apk --no-cache add ca-certificates

COPY --from=golang:1.25-alpine /usr/local/go /usr/local/go
ENV PATH=/usr/local/go/bin:$PATH

WORKDIR /app

# Create Go source file with proper routing
RUN printf '%s\n' \
'package main' \
'' \
'import (' \
'	"encoding/json"' \
'	"fmt"' \
'	"log"' \
'	"net/http"' \
'	"strings"' \
')' \
'' \
'func jsonResp(w http.ResponseWriter, code int, data any) {' \
'	w.Header().Set("Content-Type", "application/json")' \
'	w.WriteHeader(code)' \
'	enc := json.NewEncoder(w)' \
'	enc.SetIndent("", "  ")' \
'	enc.Encode(data)' \
'}' \
'' \
'func main() {' \
'	http.HandleFunc("/v1/health", func(w http.ResponseWriter, r *http.Request) {' \
'		jsonResp(w, 200, map[string]any{' \
'			"status":  "ok",' \
'			"service": "dify-mock",' \
'			"version": "1.0.0",' \
'		})' \
'	})' \
'' \
'	// Unified /v1/datasets handler' \
'	http.HandleFunc("/v1/datasets/", func(w http.ResponseWriter, r *http.Request) {' \
'		path := strings.TrimPrefix(r.URL.Path, "/v1/datasets/")' \
'		switch {' \
'		case path == "":' \
'			// List datasets' \
'			if r.Method == "GET" {' \
'				jsonResp(w, 200, map[string]any{' \
'					"code": 0,' \
'					"data": []map[string]any{{' \
'						"id":   "dataset-1",' \
'						"name": "Mock Dataset 1",' \
'					}},' \
'				})' \
'				return' \
'			}' \
'		case strings.HasSuffix(path, "/document/create-by-text"):' \
'			// Document upload' \
'			if r.Method == "POST" {' \
'				jsonResp(w, 200, map[string]any{' \
'					"code":    0,' \
'					"message": "Document created successfully (mock)",' \
'					"data":    map[string]string{"id": "mock-doc-1"},' \
'				})' \
'				return' \
'			}' \
'		default:' \
'			// Dataset detail' \
'			if r.Method == "GET" {' \
'				jsonResp(w, 200, map[string]any{' \
'					"code": 0,' \
'					"data": map[string]any{' \
'						"id":   path,' \
'						"name": "Mock Dataset " + path,' \
'					},' \
'				})' \
'				return' \
'			}' \
'		}' \
'		jsonResp(w, 404, map[string]string{"error": "not found"})' \
'	})' \
'' \
'	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {' \
'		jsonResp(w, 200, map[string]any{' \
'			"status":  "ok",' \
'			"message": "mock dify service",' \
'		})' \
'	})' \
'' \
'	fmt.Println("Mock Dify service listening on :8001")' \
'	log.Fatal(http.ListenAndServe(":8001", nil))' \
'}' > main.go

RUN go mod init dify-mock && \
    go get -d && \
    CGO_ENABLED=0 go build -o dify-mock .

EXPOSE 8001

CMD ["./dify-mock"]
