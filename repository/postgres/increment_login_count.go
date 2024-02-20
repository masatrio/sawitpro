package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/sawitpro/UserService/common"
)

func (c *Client) IncrementLoginCount(ctx context.Context, userID int64) error {
	query := `
		UPDATE users
		SET login_count = login_count + 1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $1
	`

	tx, ok := ctx.Value(common.TX_KEY).(*sql.Tx)

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

	_, err = stmt.ExecContext(ctx, userID)

	return err
}
