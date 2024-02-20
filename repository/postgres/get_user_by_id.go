package postgres

import (
	"context"
	"database/sql"

	_ "github.com/lib/pq"

	"github.com/sawitpro/UserService/repository"
)

func (c *Client) GetUserByID(ctx context.Context, userID int64) (*repository.User, error) {
	query := `SELECT id, full_name, hashed_password, phone, login_count, created_at, updated_at FROM users WHERE id = $1 LIMIT 1`

	stmt, err := c.DB.PrepareContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var user repository.User
	err = stmt.QueryRowContext(ctx, userID).Scan(
		&user.ID,
		&user.FullName,
		&user.HashedPassword,
		&user.Phone,
		&user.LoginCount,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}
