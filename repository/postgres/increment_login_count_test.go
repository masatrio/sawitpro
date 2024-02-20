package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestIncrementLoginCount(t *testing.T) {
	testCases := []struct {
		name           string
		userID         int64
		expectedError  error
		transactionCtx bool
	}{
		{
			name:           "Successful Increment",
			userID:         1,
			expectedError:  nil,
			transactionCtx: false,
		},
		{
			name:           "Error Executing Query",
			userID:         2,
			expectedError:  errors.New("some error"),
			transactionCtx: false,
		},
		{
			name:           "Successful Increment with Transaction",
			userID:         3,
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
			case "Successful Increment":
				mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE users SET login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`)).WithArgs(tc.userID).
					WillReturnResult(sqlmock.NewResult(0, 1))
			case "Error Executing Query":
				mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE users SET login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`)).WillReturnError(errors.New("some error"))
			case "Successful Increment with Transaction":
				mock.ExpectBegin()
				mock.ExpectPrepare(regexp.QuoteMeta(`UPDATE users SET login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`))
				mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP WHERE id = $1`)).WithArgs(tc.userID).
					WillReturnResult(sqlmock.NewResult(0, 1))
				mock.ExpectCommit()
			}

			if tc.transactionCtx {
				err = repo.ExecTransaction(ctx, func(ctx context.Context) error {
					return repo.IncrementLoginCount(ctx, tc.userID)
				})
			} else {
				err = repo.IncrementLoginCount(ctx, tc.userID)
			}

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
