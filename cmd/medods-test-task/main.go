package main

import (
	"jwt/internal/pkg/app"
)

func main() {
	ap, _ := app.New()
	ap.Run()
}

// func init() {
// 	err := godotenv.Load()
// 	if err != nil {
// 		log.Fatal("Error loading .env file")
// 	}

// 	service.SetJwtKey(os.Getenv("secret_key"))
// }
