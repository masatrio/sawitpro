// This file contains types that are used in the repository layer.
package repository

import "time"

type User struct {
	ID             int64
	FullName       string
	HashedPassword string
	Phone          string
	LoginCount     int64
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type UserActivityLog struct {
	ID           int64
	UserID       int64
	ActivityType string
	CreatedAt    time.Time
}
