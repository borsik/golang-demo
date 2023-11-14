package handler

import (
	"github.com/go-chi/chi/middleware"
	log "github.com/sirupsen/logrus"
	"net"
	"net/http"
)

// LoggerWithLevel logger middleware implementation for chi
func LoggerWithLevel(level log.Level) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			reqID := middleware.GetReqID(r.Context())
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)
			defer func() {
				remoteIP, _, err := net.SplitHostPort(r.RemoteAddr)
				if err != nil {
					remoteIP = r.RemoteAddr
				}
				scheme := "http"
				if r.TLS != nil {
					scheme = "https"
				}
				fields := log.Fields{
					"status_code": ww.Status(),
					"remote_ip":   remoteIP,
					"proto":       r.Proto,
					"method":      r.Method,
				}
				if len(reqID) > 0 {
					fields["request_id"] = reqID
				}
				log.WithFields(fields).Logf(level, "%s://%s%s", scheme, r.Host, r.RequestURI)
			}()

			h.ServeHTTP(ww, r)
		}

		return http.HandlerFunc(fn)
	}
}
