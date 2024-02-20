package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestUpdateUser(t *testing.T) {
	testCases := []struct {
		name           string
		userID         int64
		fullName       string
		phoneNumber    string
		expectedError  error
		transactionCtx bool
	}{
		{
			name:           "Successful Update without Transaction",
			userID:         1,
			fullName:       "maulana aji satrio",
			phoneNumber:    "+628232482440",
			expectedError:  nil,
			transactionCtx: false,
		},
		{
			name:           "Error Executing Query",
			userID:         2,
			fullName:       "maulana",
			phoneNumber:    "+628232482441",
			expectedError:  errors.New("some error"),
			transactionCtx: false,
		},
		{
			name:           "Successful Update with Transaction",
			userID:         3,
			fullName:       "maulana satrio",
			phoneNumber:    "+628232482443",
			expectedError:  nil,
			transactionCtx: true,
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

			ctx := context.Background()

			switch tc.name {
			case "Successful Update without Transaction":
				mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE users SET full_name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET full_name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`)).
					WithArgs(tc.fullName, tc.phoneNumber, tc.userID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			case "Error Executing Query":
				mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE users SET full_name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET full_name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`)).
					WithArgs(tc.fullName, tc.phoneNumber, tc.userID).
					WillReturnError(errors.New("some error"))
			case "Successful Update with Transaction":
				mock.ExpectBegin()
				mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE users SET full_name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET full_name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`)).
					WithArgs(tc.fullName, tc.phoneNumber, tc.userID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}

			if tc.transactionCtx {
				err = repo.ExecTransaction(ctx, func(ctx context.Context) error {
					return repo.UpdateUser(ctx, tc.userID, tc.fullName, tc.phoneNumber)
				})
			} else {
				err = repo.UpdateUser(ctx, tc.userID, tc.fullName, tc.phoneNumber)
			}

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
