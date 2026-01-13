package middlewares

import (
	"log"
	"net/http"
	"time"
)

func LoggerMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		log.Printf("➡️ %s %s", r.Method, r.RequestURI)

		// for key, values := range r.Header {
		// 	for _, value := range values {
		// 		fmt.Printf("%s: %s\n", key, value)
		// 	}
		// }

		next.ServeHTTP(w, r)

		log.Printf("✅ %s completed in %v", r.RequestURI, time.Since(start))
	})
}
