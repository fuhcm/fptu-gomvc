package middlewares

import "github.com/rs/cors"

// CORSMiddleware ...
func CORSMiddleware() *cors.Cors {
	return cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000", "http://localhost:3001", "https://fptu.tech", "http://fptu.tech", "https://staging.fptu.tech", "http://staging.fptu.tech", "http://news.fptu.tech", "https://news.fptu.tech"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
	})
}
