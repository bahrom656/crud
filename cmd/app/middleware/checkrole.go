package middleware

import (
	"context"
	"net/http"
)

type HasAnyRoleFunc func(ctx context.Context, roles ...string) bool

func CheckRole(roleFunc HasAnyRoleFunc, roles ...string) func(http.Handler) http.Handler {
	return func(handler http.Handler) http.Handler{
		return http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request){
			if !roleFunc(request.Context(), roles...) {
				http.Error(writer, http.StatusText(http.StatusForbidden), http.StatusForbidden)
				return
			}

			handler.ServeHTTP(writer, request)
		})
	}
}
