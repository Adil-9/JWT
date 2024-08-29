package service

import (
	"encoding/base64"
	"fmt"
	errorlog "jwt/internal/app/errorLog"
	"jwt/internal/structs"
	"time"

	"net/smtp"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

var jwtKey []byte

func SetJwtKey(key string) {
	jwtKey = []byte(key)
}

type Database interface {
	PutInDB(hash []byte, ip string) error
	GetDataFromDB(ip string) (structs.RefreshToken, error)
}

type Service struct {
	DB Database
}

func New(db Database) *Service {
	return &Service{DB: db}
}

func (s *Service) GetjwtKey() []byte {
	return jwtKey
}

func (s *Service) GetClientIP(c *gin.Context) string {
	IPAddress := c.Request.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = c.Request.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = c.Request.RemoteAddr
	}
	return IPAddress
}

func (s *Service) GetIPData(ip string) (structs.RefreshToken, error) {
	return s.DB.GetDataFromDB(ip)
}

func (s *Service) SendWarningEmail(userID string) {
	  from := "some@mail.com"
	  password := "somepassword"
	
	  to := []string{
		"tosomeemail.com",
	  }
	
	  smtpHost := "smtp.gmail.com"
	  smtpPort := "587"
	
	  message := []byte("This is a test email message.")
	  
	  auth := smtp.PlainAuth("", from, password, smtpHost)
	  
	  err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, to, message)
	  if err != nil {
		errorlog.ErrorPrint("can not send email", err)
		return
	  }
	  fmt.Println("Email Sent Successfully!")
}

func (s *Service) GenerateTokens(ip string) (string, string, error) {
	expirationTime := time.Now().Add(17 * time.Minute)

	data := &structs.Data{
		IP: ip,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS512, data)
	accessTokenSigned, err := accessToken.SignedString(jwtKey)
	if err != nil {
		return "", "", err
	}

	refreshTokenData := fmt.Sprintf("%s:%s", ip, uuid.NewString())
	refreshToken := base64.StdEncoding.EncodeToString([]byte(refreshTokenData))

	hashedToken, err := bcrypt.GenerateFromPassword([]byte(refreshToken), bcrypt.DefaultCost)
	if err != nil {
		return "", "", err
	}

	err = s.DB.PutInDB(hashedToken, ip)
	if err != nil {
		return "", "", err
	}

	return accessTokenSigned, refreshToken, nil
}
