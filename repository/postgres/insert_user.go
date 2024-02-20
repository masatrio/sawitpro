package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/sawitpro/UserService/common"
	"github.com/sawitpro/UserService/repository"
)

func (c *Client) InsertUser(ctx context.Context, user *repository.User) (*repository.User, error) {
	// Define the query to insert a new user
	query := `
		INSERT INTO users (full_name, hashed_password, phone)
		VALUES ($1, $2, $3)
		RETURNING id, full_name, hashed_password, phone, login_count, created_at, updated_at
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
		return nil, err
	}
	defer stmt.Close()

	// Execute the prepared statement within the transaction or directly
	var row *sql.Row
	if tx != nil {
		row = stmt.QueryRowContext(ctx, user.FullName, user.HashedPassword, user.Phone)
	} else {
		row = stmt.QueryRowContext(ctx, user.FullName, user.HashedPassword, user.Phone)
	}

	// Scan the result into a user object
	insertedUser := &repository.User{}
	err = row.Scan(
		&insertedUser.ID,
		&insertedUser.FullName,
		&insertedUser.HashedPassword,
		&insertedUser.Phone,
		&insertedUser.LoginCount,
		&insertedUser.CreatedAt,
		&insertedUser.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return insertedUser, nil
}
