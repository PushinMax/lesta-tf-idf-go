package service

import "github.com/PushinMax/lesta-tf-idf-go/internal/repository"

type UserService struct {
	repos *repository.Repository
}

func newUserApi(repos *repository.Repository) *UserService {
	return &UserService{
		repos: repos,
	}
}

func (s *UserService) ChangePassword(id, password string) error {
	return s.repos.ChangePassword(id, password)
}
