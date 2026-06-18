package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/hugaojanuario/sentinel/internal/models"
	"github.com/hugaojanuario/sentinel/internal/repository"
	"github.com/hugaojanuario/sentinel/pkg/config"
	"golang.org/x/crypto/bcrypt"
)

var secret = []byte(config.LoadDotEnv().JWT_SECRET)
var ErrInvalidCredentials = errors.New("invalid credentials")

type Claims struct {
	UserId string `json:"userId"`
	jwt.RegisteredClaims
}

type Service struct {
	r *repository.Repository
}

func NewService(r *repository.Repository) *Service {
	return &Service{r: r}
}

func (s *Service) CreateUser(req models.CreateUserRequest) (*models.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("erro ao gerar o hash da senha: %w", err)
	}

	user, err := s.r.CreateUser(req, string(hashedPassword))
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (s *Service) Login(req models.LoginRequest) (*models.LoginResponse, error) {
	existing, err := s.r.FindUserByEMail(req.Email)
	if err != nil {
		return nil, fmt.Errorf("erro ao encontrar o usuario pelo email: %w", err)
	}
	if existing == nil {
		return nil, ErrInvalidCredentials
	}

	err = bcrypt.CompareHashAndPassword([]byte(existing.PasswordHash), []byte(req.Password))
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	claims := Claims{
		UserId: existing.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenStr, err := token.SignedString(secret)
	if err != nil {
		return nil, fmt.Errorf("erro ao assinar token: %w", err)
	}

	return &models.LoginResponse{Token: tokenStr}, nil
}
