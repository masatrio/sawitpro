package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/sawitpro/UserService/common"
)

func (c *Client) UpdateUser(ctx context.Context, userID int64, fullName, phoneNumber string) error {

	query := `UPDATE users SET full_name = $1, phone = $2, updated_at = CURRENT_TIMESTAMP WHERE id = $3`

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

	// Execute the statement with the updated values
	_, err = stmt.ExecContext(ctx, fullName, phoneNumber, userID)

	return err
}
