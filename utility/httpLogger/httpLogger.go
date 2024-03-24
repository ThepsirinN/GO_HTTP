package httpLogger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go_http_barko/utility/logger"
	"io"
	"net/http"
	"time"

	"go.uber.org/zap"
)

type httpPayload struct {
	RequestMethod string
	RequestURL    string
	UserAgent     string
	RemoteIP      string
	Referrer      string
	Protocol      string
	Status        string
	ResponseSize  string
	Latency       string
}

func HttpLoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var httpPayload httpPayload

		httpPayload.RequestMethod = r.Method
		httpPayload.UserAgent = r.UserAgent()
		httpPayload.RemoteIP = r.RemoteAddr
		httpPayload.Protocol = r.Proto
		httpPayload.Referrer = r.Referer()
		httpPayload.RequestURL = r.RequestURI

		logger.Info(r.Context(), fmt.Sprintf("Receive Http request Method=%s Path=%s", httpPayload.RequestMethod, httpPayload.RequestURL),
			zap.Any("UserAgent", httpPayload.UserAgent),
			zap.Any("RemoteIP", httpPayload.RemoteIP),
			zap.Any("Protocol", httpPayload.Protocol),
			zap.Any("Referrer", httpPayload.Referrer),
			zap.Any("body", r.Body),
		)

		start := time.Now()

		// get response by add new buffer writer
		buff := new(bytes.Buffer)
		// writer := io.MultiWriter(buff) // write on buff only
		writer := io.MultiWriter(buff, w) // for writing both io.Writer and http.ResponseWriter
		writers := &respWriter{Writer: writer, ResponseWriter: w}

		next.ServeHTTP(writers, r)
		latency := time.Since(start)
		httpPayload.Latency = latency.String()
		httpPayload.Status = w.Header().Get("Status")

		var respHeader []map[string][]string
		for k, v := range w.Header() {
			respHeader = append(respHeader, map[string][]string{k: v})
		}

		var respBody []map[string]interface{}
		if w.Header().Get("Content-Type") == "application/json" {
			if buff.Len() != 0 {
				json.Unmarshal(buff.Bytes(), &respBody)
			}
		}

		status := "Status Not found!"
		if len(w.Header().Get("Status")) != 0 {
			status = string(w.Header().Get("Status"))
		}

		// httpPayload.ResponseSize
		logger.Info(r.Context(), fmt.Sprintf("Receive Http response Method=%s Path=%s", httpPayload.RequestMethod, httpPayload.RequestURL),
			// zap.Any("header", r.Response.Header),
			zap.Any("header", respHeader),
			zap.Any("body", respBody),
			zap.Any("status", status),
			zap.Any("latency", httpPayload.Latency),
		)

	})
}

type respWriter struct {
	io.Writer
	http.ResponseWriter
}

func (w *respWriter) WriteHeader(statusCode int) {
	w.ResponseWriter.WriteHeader(statusCode)
}

func (w *respWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}
