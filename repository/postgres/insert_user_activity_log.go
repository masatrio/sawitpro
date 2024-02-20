package postgres

import (
	"context"
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"

	"github.com/sawitpro/UserService/common"
)

func (c *Client) InsertUserActivityLog(ctx context.Context, userID int64, activityType string) error {
	// Define the query to insert user activity log
	query := `
		INSERT INTO user_activity_logs (user_id, activity_type)
		VALUES ($1, $2)
		RETURNING id
	`
	tx, ok := ctx.Value(common.TX_KEY).(*sql.Tx)
	// Prepare the statement
	var stmt *sql.Stmt
	var err error
	if ok && tx != nil {
		stmt, err = tx.PrepareContext(ctx, query)
	} else {
		stmt, err = c.DB.PrepareContext(ctx, query)
	}
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Execute the prepared statement within the transaction or directly
	var result sql.Result
	result, err = stmt.ExecContext(ctx, userID, activityType)
	if err != nil {
		return err
	}

	// Check if the insert was successful
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}
