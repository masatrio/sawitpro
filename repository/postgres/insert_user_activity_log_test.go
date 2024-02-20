package postgres

import (
	"context"
	"errors"
	"regexp"
	"testing"

	sqlmock "github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestInsertUserActivityLog(t *testing.T) {
	testCases := []struct {
		name           string
		userID         int64
		activityType   string
		expectedError  error
		transactionCtx bool
	}{
		{
			name:           "Successful Insert",
			userID:         1,
			activityType:   "login",
			expectedError:  nil,
			transactionCtx: false,
		},
		{
			name:           "Error Executing Query",
			userID:         2,
			activityType:   "login",
			expectedError:  errors.New("some error"),
			transactionCtx: false,
		},
		{
			name:           "No Rows Affected",
			userID:         3,
			activityType:   "login",
			expectedError:  errors.New("no rows affected"),
			transactionCtx: false,
		},
		{
			name:           "Successful Insert with Transaction",
			userID:         4,
			activityType:   "login",
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
			case "Successful Insert":
				mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO user_activity_logs (user_id, activity_type) VALUES ($1, $2) RETURNING id`))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO user_activity_logs (user_id, activity_type) VALUES ($1, $2) RETURNING id`)).
					WithArgs(tc.userID, tc.activityType).
					WillReturnResult(sqlmock.NewResult(1, 1))
			case "Error Executing Query":
				mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO user_activity_logs (user_id, activity_type) VALUES ($1, $2) RETURNING id`))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO user_activity_logs (user_id, activity_type) VALUES ($1, $2) RETURNING id`)).
					WithArgs(tc.userID, tc.activityType).
					WillReturnError(errors.New("some error"))
			case "No Rows Affected":
				mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO user_activity_logs (user_id, activity_type) VALUES ($1, $2) RETURNING id`))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO user_activity_logs (user_id, activity_type) VALUES ($1, $2) RETURNING id`)).
					WithArgs(tc.userID, tc.activityType).
					WillReturnResult(sqlmock.NewResult(0, 0))
			case "Successful Insert with Transaction":
				mock.ExpectBegin()
				mock.ExpectPrepare(regexp.QuoteMeta(`INSERT INTO user_activity_logs (user_id, activity_type) VALUES ($1, $2) RETURNING id`))
				mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO user_activity_logs (user_id, activity_type) VALUES ($1, $2) RETURNING id`)).
					WithArgs(tc.userID, tc.activityType).
					WillReturnResult(sqlmock.NewResult(1, 1))
				mock.ExpectCommit()
			}

			if tc.transactionCtx {
				err = repo.ExecTransaction(ctx, func(ctx context.Context) error {
					err := repo.InsertUserActivityLog(ctx, tc.userID, tc.activityType)
					return err
				})
			} else {
				err = repo.InsertUserActivityLog(ctx, tc.userID, tc.activityType)
			}

			assert.Equal(t, tc.expectedError, err)
		})
	}
}
