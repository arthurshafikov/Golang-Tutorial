package internalhttp

import (
	"fmt"
	"net/http"
	"time"
)

// todo make method?
func loggingMiddleware(next http.Handler, server *Server) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		server.Logger.Info(
			fmt.Sprintf("%v [%v] %v %v %v",
				r.RemoteAddr,
				time.Now().Format("02/01/2006:15:04:05 MST"),
				r.Method,
				r.RequestURI,
				r.UserAgent(),
			),
		)
		next.ServeHTTP(w, r)
	})
}

// func jsonDecodeMiddleware(next http.Handler, server *Server) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		var decoded interface{}
// 		decoder := json.NewDecoder(r.Body)
// 		if err := decoder.Decode(&decoded); err != nil {
// 			panic(err)
// 		}
// 		buf := strings.NewReader(decoded)
// 		r.Body = io.NopCloser(buf)
// 		next.ServeHTTP(w, r)
// 	})
// }
