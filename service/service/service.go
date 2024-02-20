package service

import (
	"github.com/sawitpro/UserService/repository"
	"github.com/sawitpro/UserService/service"
)

type Service struct {
	Repository repository.RepositoryInterface
}

type ServiceOpts struct {
	Repository repository.RepositoryInterface
}

func NewService(opts ServiceOpts) service.ServiceInterface {
	return &Service{
		Repository: opts.Repository,
	}
}
