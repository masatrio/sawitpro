// This file contains the interfaces for the repository layer.
// The repository layer is responsible for interacting with the database.
// For testing purpose we will generate mock implementations of these
// interfaces using mockgen. See the Makefile for more information.
package repository

import (
	"context"
)

type RepositoryInterface interface {
	ExecTransaction(ctx context.Context, fn func(context.Context) error) error
	InsertUser(ctx context.Context, user *User) (*User, error)
	GetUserByPhone(ctx context.Context, phone string) (*User, error)
	GetUserByID(ctx context.Context, userID int64) (*User, error)
	IncrementLoginCount(ctx context.Context, userID int64) error
	InsertUserActivityLog(ctx context.Context, userID int64, activityType string) error
	UpdateUser(ctx context.Context, userID int64, fullName, phoneNumber string) error
}
