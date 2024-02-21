package service

import (
	"github.com/sawitpro/UserService/helper/hasher"
	"github.com/sawitpro/UserService/repository"
	"github.com/sawitpro/UserService/service"
)

type Service struct {
	Repository repository.RepositoryInterface
	Hasher     hasher.PasswordHasher
}

type ServiceOpts struct {
	Repository repository.RepositoryInterface
	Hasher     hasher.PasswordHasher
}

func NewService(opts ServiceOpts) service.ServiceInterface {
	return &Service{
		Repository: opts.Repository,
		Hasher:     opts.Hasher,
	}
}
