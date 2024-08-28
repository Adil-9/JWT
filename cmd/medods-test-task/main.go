package main

import (
	"jwt/internal/app/service"
	"jwt/internal/pkg/app"
	"log"
	"os"

	"github.com/lpernett/godotenv"
)

// var jwtKey []byte

// type Data struct {
// 	IP string `json:"ip"`
// 	jwt.RegisteredClaims
// }

// type RefreshToken struct {
// 	Hash   []byte
// 	UserID string
// 	IP     string
// }

// var RefreshTokensDB = map[string]structs.RefreshToken{}

// func sendWarningEmail(userID string) {
// 	fmt.Printf("Sending warning email to user %s due to IP address change.\n", userID)
// }

// func GenerateTokens(ip string) (string, string, error) {
// 	expirationTime := time.Now().Add(17 * time.Minute)

// 	data := &Data{
// 		IP: ip,
// 		RegisteredClaims: jwt.RegisteredClaims{
// 			ExpiresAt: jwt.NewNumericDate(expirationTime),
// 		},
// 	}

// 	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, data)
// 	accessTokenSigned, err := accessToken.SignedString(jwtKey)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	refreshTokenData := fmt.Sprintf("%s:%s", ip, uuid.NewString())
// 	refreshToken := base64.StdEncoding.EncodeToString([]byte(refreshTokenData))

// 	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
// 	if err != nil {
// 		return "", "", err
// 	}

// 	refreshTokensDB[ip] = RefreshToken{Hash: hashedToken, IP: ip}

// 	return accessTokenSigned, refreshToken, nil
// }

// func Access(c *gin.Context) {
// 	clientIP := getClientIP(c)

// 	accessToken, refreshToken, err := GenerateTokens(clientIP)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "can not create tokens"})
// 		return
// 	}

// 	c.SetCookie(
// 		"access_token",
// 		accessToken,
// 		1200,
// 		"/",
// 		"",
// 		true,
// 		true,
// 	)

// 	// Set the refresh token as an HttpOnly cookie
// 	c.SetCookie(
// 		"refresh_token",
// 		refreshToken,
// 		86400,
// 		"/",
// 		"",
// 		true,
// 		true,
// 	)

// 	c.JSON(http.StatusOK, gin.H{
// 		"access_token":  accessToken,
// 		"refresh_token": refreshToken,
// 	})
// }

// func Refresh(c *gin.Context) {
// 	accessToken, err := c.Cookie("access_token")
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found in cookies"})
// 		return
// 	}

// 	refreshToken, err := c.Cookie("refresh_token")
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found in cookies"})
// 		return
// 	}

// 	data := &Data{}
// 	tkn, err := jwt.ParseWithClaims(accessToken, data, func(token *jwt.Token) (interface{}, error) {
// 		return jwtKey, nil
// 	})
// 	if err != nil || !tkn.Valid {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Access token"})
// 		return
// 	}

// 	storedToken, exists := refreshTokensDB[data.IP]
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Refresh token"})
// 		return
// 	}

// 	err = bcrypt.CompareHashAndPassword(storedToken.Hash, []byte(refreshToken))
// 	if err != nil {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Refresh token"})
// 		return
// 	}

// 	clientIP := getClientIP(c)
// 	if storedToken.IP != clientIP {
// 		sendWarningEmail(data.IP)
// 	}

// 	newAccessToken, newRefreshToken, err := GenerateTokens(data.IP)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not refresh tokens"})
// 		return
// 	}

// 	c.SetCookie(
// 		"access_token",
// 		newAccessToken,
// 		1200,
// 		"/",
// 		"",
// 		true,
// 		true,
// 	)

// 	c.SetCookie(
// 		"refresh_token",
// 		newRefreshToken,
// 		86400,
// 		"/",
// 		"",
// 		true,
// 		true,
// 	)

// 	c.JSON(http.StatusOK, gin.H{
// 		"access_token":  newAccessToken,
// 		"refresh_token": newRefreshToken,
// 	})
// }

// func getClientIP(c *gin.Context) string {
// 	IPAddress := c.Request.Header.Get("X-Real-Ip")
// 	if IPAddress == "" {
// 		IPAddress = c.Request.Header.Get("X-Forwarded-For")
// 	}
// 	if IPAddress == "" {
// 		IPAddress = c.Request.RemoteAddr
// 	}
// 	return IPAddress
// }

func main() {
	ap, _ := app.New()
	ap.Run()
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	service.SetJwtKey(os.Getenv("secret_key"))
}
