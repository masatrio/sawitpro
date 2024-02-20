package handler

import (
	"github.com/sawitpro/UserService/generated"
	"github.com/sawitpro/UserService/service"
)

type Server struct {
	Service service.ServiceInterface
}

type ServerOpts struct {
	Service service.ServiceInterface
}

func NewServer(opts ServerOpts) generated.ServerInterface {
	return &Server{
		Service: opts.Service,
	}
}
