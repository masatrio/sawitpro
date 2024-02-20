package postgres

import (
	"context"
	"database/sql"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/sawitpro/UserService/repository"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByPhone(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name          string
		phoneNumber   string
		expectedUser  *repository.User
		expectedError error
	}{
		{
			name:        "User Exists",
			phoneNumber: "+628232482440",
			expectedUser: &repository.User{
				ID:             1,
				FullName:       "maulana aji satrio",
				HashedPassword: "Maulana1996@",
				Phone:          "+628232482440",
				LoginCount:     1,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			expectedError: nil,
		},
		{
			name:          "User Not Found",
			phoneNumber:   "+628232482441",
			expectedUser:  nil,
			expectedError: nil,
		},
		{
			name:          "Error Executing Query",
			phoneNumber:   "+628232482443",
			expectedUser:  nil,
			expectedError: errors.New("some error"),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock database: %v", err)
			}
			defer mockDB.Close()

			repo := &Client{
				DB: mockDB,
			}

			mock.ExpectPrepare(regexp.QuoteMeta(`SELECT id, full_name, hashed_password, phone, login_count, created_at, updated_at FROM users WHERE phone = $1 LIMIT 1`))
			switch tc.name {
			case "User Exists":
				rows := sqlmock.NewRows([]string{"id", "full_name", "hashed_password", "phone", "login_count", "created_at", "updated_at"}).
					AddRow(1, "maulana aji satrio", "Maulana1996@", "+628232482440", 1, now, now)
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, full_name, hashed_password, phone, login_count, created_at, updated_at FROM users WHERE phone = $1 LIMIT 1`)).WillReturnRows(rows)
			case "User Not Found":
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, full_name, hashed_password, phone, login_count, created_at, updated_at FROM users WHERE phone = $1 LIMIT 1`)).WillReturnError(sql.ErrNoRows)
			case "Error Executing Query":
				mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, full_name, hashed_password, phone, login_count, created_at, updated_at FROM users WHERE phone = $1 LIMIT 1`)).WillReturnError(errors.New("some error"))
			}

			user, err := repo.GetUserByPhone(context.Background(), tc.phoneNumber)

			assert.Equal(t, tc.expectedUser, user)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
