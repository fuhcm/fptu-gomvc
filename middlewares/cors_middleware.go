package middlewares

import "github.com/rs/cors"

// CORSMiddleware ...
func CORSMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
	})
}