package middleware

import (
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Header Constants
const (
	HeaderAllowOrigin  = "Access-Control-Allow-Origin"
	HeaderAllowMethods = "Access-Control-Allow-Methods"
	HeaderAllowHeaders = "Access-Control-Allow-Headers"
	AllowedMethods     = "GET, OPTIONS"
	AllowedHeaders     = "Origin, Referer, Accept, Accept-Encoding, Accept-Language, x-requested-with, Content-Type, Content-Length, Authorization"
)

func CORSPolicies(allMethods []string, allowedOrigins []string) mux.MiddlewareFunc {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
			var allowedOrigin = ""
			var ref = r.Header.Get("Origin")
			for _, v := range allowedOrigins {
				if strings.Contains(ref, v) {
					allowedOrigin = ref
				}
			}
			for _, v := range allMethods {
				if v == http.MethodOptions {
					rw.Header().Add(HeaderAllowOrigin, allowedOrigin)
					if r.Method == http.MethodOptions {
						rw.Header().Add(HeaderAllowHeaders, AllowedHeaders)
						rw.Header().Add(HeaderAllowMethods, strings.Join(allMethods, ","))
					}
				}
			}
			h.ServeHTTP(rw, r)
		})
	}
}
