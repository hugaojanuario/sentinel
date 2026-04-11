package services

import (
	"fmt"

	"github.com/hugaojanuario/sentinel/internal/models"
	"github.com/hugaojanuario/sentinel/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

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

func (s *Service) FindUserByEmail(email string) (models.User, error) {
	user, err := s.r.FindUserByEMail(email)
	if err !+ nil{
		return nil, err
	}


}
