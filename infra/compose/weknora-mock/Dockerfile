# Mock WeKnora service for CI/integration testing
# Provides minimal endpoints needed for Servify integration

FROM golang:1.23-alpine AS builder
WORKDIR /app
RUN echo 'package main; import ("net/http","encoding/json"); func main(){http.HandleFunc("/api/v1/health",func(w http.ResponseWriter,r*http.Request){w.Header().Set("Content-Type","application/json");json.NewEncoder(w).Encode(map[string]interface{}{"status":"ok","service":"weknora-mock","version":"1.0.0"})});http.HandleFunc("/",func(w http.ResponseWriter,r*http.Request){w.Header().Set("Content-Type","application/json");json.NewEncoder(w).Encode(map[string]interface{}{"status":"ok","message":"mock weknora service"})});http.ListenAndServe(":9000",nil)}' > main.go
RUN CGO_ENABLED=0 go build -o weknora-mock .

FROM alpine:latest
COPY --from=builder /app/weknora-mock /usr/local/bin/
EXPOSE 9000 9001
CMD ["weknora-mock"]
