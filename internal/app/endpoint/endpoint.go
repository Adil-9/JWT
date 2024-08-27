package endpoint

import (
	errorlog "jwt/internal/app/errorLog"
	"jwt/internal/structs"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	GetClientIP(*gin.Context) string
	GenerateTokens(ip string) (string, string, error)
	SendWarningEmail(userID string)
	GetIPData(ip string) (structs.RefreshToken, error)
	GetjwtKey() []byte
}

type EndPoint struct {
	s Service
}

func New(service Service) *EndPoint {
	return &EndPoint{s: service}
}

func (e *EndPoint) Access(c *gin.Context) {
	clientIP := e.s.GetClientIP(c)

	accessToken, refreshToken, err := e.s.GenerateTokens(clientIP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "can not create tokens"})
		errorlog.ErrorPrint("can not create tokens", err)
		return
	}

	c.SetCookie(
		"access_token",
		accessToken,
		1200,
		"/",
		"",
		true,
		true,
	)

	// Set the refresh token as an HttpOnly cookie
	c.SetCookie(
		"refresh_token",
		refreshToken,
		86400,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func (e *EndPoint) Refresh(c *gin.Context) {
	accessToken, err := c.Cookie("access_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Access token not found in cookies"})
		errorlog.ErrorPrint("access token not found in cookies", err)
		return
	}

	refreshToken, err := c.Cookie("refresh_token")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found in cookies"})
		errorlog.ErrorPrint("refresh token not found in cookies", err)
		return
	}

	data := &structs.Data{}
	tkn, err := jwt.ParseWithClaims(accessToken, data, func(token *jwt.Token) (interface{}, error) {
		return e.s.GetjwtKey(), nil
	})
	if err != nil || !tkn.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Access token"})
		errorlog.ErrorPrint("invalid access token", err)
		return
	}

	// storedToken, exists := refreshTokensDB[data.IP]
	storedToken, err := e.s.GetIPData(data.IP)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Access token"})
		errorlog.ErrorPrint("invalid token from database", err)
		return
	}

	// if !exists {
	// 	c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Refresh token"})
	// 	return
	// }

	err = bcrypt.CompareHashAndPassword(storedToken.Hash, []byte(refreshToken))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid Refresh token"})
		errorlog.ErrorPrint("invalid refresh token", err)
		return
	}

	clientIP := e.s.GetClientIP(c)
	if storedToken.IP != clientIP {
		e.s.SendWarningEmail(data.IP)
	}

	newAccessToken, newRefreshToken, err := e.s.GenerateTokens(data.IP)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not refresh tokens"})
		errorlog.ErrorPrint("could not refresh tokens", err)
		return
	}

	c.SetCookie(
		"access_token",
		newAccessToken,
		1200,
		"/",
		"",
		true,
		true,
	)

	c.SetCookie(
		"refresh_token",
		newRefreshToken,
		86400,
		"/",
		"",
		true,
		true,
	)

	c.JSON(http.StatusOK, gin.H{
		"access_token":  newAccessToken,
		"refresh_token": newRefreshToken,
	})
}
