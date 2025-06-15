package service

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/PushinMax/lesta-tf-idf-go/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/spf13/viper"
)

//"errors"



type AuthService struct {
	repos *repository.Repository
	cfg  JWTConfig
}

func newAuthApi(repos *repository.Repository) *AuthService {
	accessExpiry, _ := time.ParseDuration(viper.GetString("jwt.access_expiry"))
	refreshExpiry, _ := time.ParseDuration(viper.GetString("jwt.refresh_expiry"))
	return &AuthService{
		repos: repos,
		cfg: JWTConfig{
			AccessSecret:  os.Getenv("JWT_ACCESS_SECRET"),
			RefreshSecret: os.Getenv("JWT_REFRESH_SECRET"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
	}
}

func (s *AuthService) Login(login, password, ip string) (*TokenPairResponse, error) {
	userID, err := s.repos.Authentication(login, password)
	if err != nil {
		return nil, err
	}
	jti := uuid.New().String()
	
	accessToken, err := GenerateJWT(
		userID.String(),
		ip,
		jti,
		s.cfg.AccessSecret,
		s.cfg.AccessExpiry,
		"access",
	)
	if err != nil {
		return nil, fmt.Errorf("access token generation failed: %w", err)
	}

	refreshToken, err := GenerateJWT(
		userID.String(),
		ip,
		jti,
		s.cfg.RefreshSecret,
		s.cfg.RefreshExpiry,
		"refresh",
	)

	if err != nil {
		return nil, fmt.Errorf("refresh token generation failed: %w", err)
	}

	 if err := s.repos.SetRefreshToken(userID, refreshToken); err != nil {
		return nil, fmt.Errorf("failed to set refresh token: %w", err)
    }

	return &TokenPairResponse{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}, nil
}

func (s *AuthService) Register(login, password string) error {
	return s.repos.Register(login, password)
}

func (s *AuthService) ValidateToken(token string) (*CustomClaims, error) {
	claims, err := ValidateJWT(token, s.cfg.AccessSecret)
	if err != nil {
		return nil, err
	}
	// проверить стоп-лист
	return claims, nil
}

func (s *AuthService) RefreshToken(refreshToken, ip string) (*TokenPairResponse, error) {
 	claims, err := ValidateJWT(refreshToken, s.cfg.RefreshSecret)
	if err != nil {
		return nil, err
	}
	if claims.IP != ip {
		log.Printf("IP address mismatch: expected %s, got %s", claims.IP, ip)
  		return nil, fmt.Errorf("IP address mismatch")
 	}
	if claims.TokenType != "refresh" {
		return nil, fmt.Errorf("invalid token type: %s", claims.TokenType)
	}
	if claims.IssuedAt.Time.Add(s.cfg.RefreshExpiry).Before(time.Now()) {
  		return nil, fmt.Errorf("refresh token expired")
	}


	jti := uuid.New().String()

	
	accessToken, err := GenerateJWT(
		claims.Subject,
		ip,
		jti,
		s.cfg.AccessSecret,
		s.cfg.AccessExpiry,
		"access",
	)
	if err != nil {
		return nil, fmt.Errorf("access token generation failed: %w", err)
	}

	newRefreshToken, err := GenerateJWT(
		claims.Subject,
		ip,
		jti,
		s.cfg.RefreshSecret,
		s.cfg.RefreshExpiry,
		"refresh",
	)

	if err != nil {
		return nil, fmt.Errorf("refresh token generation failed: %w", err)
	}

	if err := s.repos.CheckAndChangeRefreshToken(uuid.MustParse(claims.Subject), refreshToken, newRefreshToken); err != nil {
  		return nil, fmt.Errorf("failed to check and change refresh token: %w", err)
 	}
	return &TokenPairResponse{
		AccessToken: accessToken,
		RefreshToken: newRefreshToken,
	}, nil
}

func (s *AuthService) Logout(userID string) error {
 	if err := s.repos.Logout(userID); err != nil {
  		return fmt.Errorf("logout failed: %w", err)		
   	}
	return nil
}

type CustomClaims struct {
	IP       string `json:"ip"`
	TokenType string `json:"type"`
	JTI      string `json:"jti"`
	jwt.RegisteredClaims
}

func ValidateJWT(tokenString, secret string) (*CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return []byte(secret), nil
	})
	if claims, ok := token.Claims.(*CustomClaims); ok {
		return claims, nil
	}
	return nil, err
}

func GenerateJWT(userID, ip, jti, secret string, expiry time.Duration, tp string) (string, error) {
	claims := jwt.MapClaims{
		"sub":  userID,
		"ip":   ip,
		"jti":  jti,
		"exp":  time.Now().Add(expiry).Unix(),
		"iat":  time.Now().Unix(),
		"type": tp,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	return token.SignedString([]byte(secret))
}

