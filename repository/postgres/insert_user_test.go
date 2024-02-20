package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/sawitpro/UserService/repository"
	"github.com/stretchr/testify/assert"
)

func TestInsertUser(t *testing.T) {
	now := time.Now()
	testCases := []struct {
		name           string
		user           *repository.User
		expectedUser   *repository.User
		expectedError  error
		transactionCtx bool
	}{
		{
			name: "Successful Insert without Transaction",
			user: &repository.User{
				FullName:       "maulana aji satrio",
				HashedPassword: "Maulana1996@",
				Phone:          "+628232482440",
			},
			expectedUser: &repository.User{
				ID:             1,
				FullName:       "maulana aji satrio",
				HashedPassword: "Maulana1996@",
				Phone:          "+628232482440",
				LoginCount:     0,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			expectedError:  nil,
			transactionCtx: false,
		},
		{
			name: "Successful Insert with Transaction",
			user: &repository.User{
				FullName:       "maulana aji",
				HashedPassword: "Maulana1997@",
				Phone:          "+628232482446",
			},
			expectedUser: &repository.User{
				ID:             1,
				FullName:       "maulana aji",
				HashedPassword: "Maulana1997@",
				Phone:          "+628232482446",
				LoginCount:     0,
				CreatedAt:      now,
				UpdatedAt:      now,
			},
			expectedError:  nil,
			transactionCtx: true,
		},
		{
			name: "Error Executing Query",
			user: &repository.User{
				FullName:       "Error User",
				HashedPassword: "Maulana1996@",
				Phone:          "+628232482448",
			},
			expectedUser:   nil,
			expectedError:  errors.New("some error"),
			transactionCtx: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a new mock database connection
			mockDB, mock, err := sqlmock.New()
			if err != nil {
				t.Fatalf("failed to create mock database: %v", err)
			}
			defer mockDB.Close()

			repo := &Client{
				DB: mockDB,
			}

			ctx := context.Background()

			switch tc.name {
			case "Successful Insert without Transaction":
				mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO users (full_name, hashed_password, phone) VALUES ($1, $2, $3) RETURNING id, full_name, hashed_password, phone, login_count, created_at, updated_at`))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (full_name, hashed_password, phone) VALUES ($1, $2, $3) RETURNING id, full_name, hashed_password, phone, login_count, created_at, updated_at`)).
					WithArgs(tc.user.FullName, tc.user.HashedPassword, tc.user.Phone).
					WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "hashed_password", "phone", "login_count", "created_at", "updated_at"}).
						AddRow(tc.expectedUser.ID, tc.expectedUser.FullName, tc.expectedUser.HashedPassword, tc.expectedUser.Phone, tc.expectedUser.LoginCount, tc.expectedUser.CreatedAt, tc.expectedUser.UpdatedAt))
			case "Successful Insert with Transaction":
				mock.ExpectBegin()
				mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO users (full_name, hashed_password, phone) VALUES ($1, $2, $3) RETURNING id, full_name, hashed_password, phone, login_count, created_at, updated_at`))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (full_name, hashed_password, phone) VALUES ($1, $2, $3) RETURNING id, full_name, hashed_password, phone, login_count, created_at, updated_at`)).
					WithArgs(tc.user.FullName, tc.user.HashedPassword, tc.user.Phone).
					WillReturnRows(sqlmock.NewRows([]string{"id", "full_name", "hashed_password", "phone", "login_count", "created_at", "updated_at"}).
						AddRow(tc.expectedUser.ID, tc.expectedUser.FullName, tc.expectedUser.HashedPassword, tc.expectedUser.Phone, tc.expectedUser.LoginCount, tc.expectedUser.CreatedAt, tc.expectedUser.UpdatedAt))
				mock.ExpectCommit()
			case "Error Executing Query":
				mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO users (full_name, hashed_password, phone) VALUES ($1, $2, $3) RETURNING id, full_name, hashed_password, phone, login_count, created_at, updated_at`))
				mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO users (full_name, hashed_password, phone) VALUES ($1, $2, $3) RETURNING id, full_name, hashed_password, phone, login_count, created_at, updated_at`)).
					WithArgs(tc.user.FullName, tc.user.HashedPassword, tc.user.Phone).
					WillReturnError(errors.New("some error"))
			}

			// Call the InsertUser method with the mocked context and user object
			var user *repository.User

			if tc.transactionCtx {
				err = repo.ExecTransaction(ctx, func(ctx context.Context) error {
					user, err = repo.InsertUser(ctx, tc.user)
					if err != nil {
						return err
					}

					return nil
				})
			} else {
				user, err = repo.InsertUser(ctx, tc.user)
			}

			// Assert the results
			assert.Equal(t, tc.expectedUser, user)
			assert.Equal(t, tc.expectedError, err)
		})
	}
}
