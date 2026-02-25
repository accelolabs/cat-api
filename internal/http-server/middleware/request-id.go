package middleware

import (
	ctxUtil "accelolabs/cat-api/internal/util/ctx"
	"crypto/rand"
	"encoding/hex"
	"net/http"
)

const requestIDHeaderKey = "X-Request-ID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestID := r.Header.Get(requestIDHeaderKey)
		if requestID == "" {
			b := make([]byte, 16)
			rand.Read(b)
			requestID = hex.EncodeToString(b)
		}

		ctx = ctxUtil.SetRequestID(ctx, requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
