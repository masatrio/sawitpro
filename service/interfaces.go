// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package service

import (
	"context"

	"github.com/sawitpro/UserService/common"
)

type ServiceInterface interface {
	Register(ctx context.Context, params RegisterParam) (*RegisterResponse, common.Error)

	Login(ctx context.Context, params LoginParam) (*LoginResponse, common.Error)

	GetProfile(ctx context.Context, userID int64) (*UserInfoResponse, common.Error)

	UpdateProfile(ctx context.Context, params UpdateProfileParam) (*UpdateProfileResponse, common.Error)
}
