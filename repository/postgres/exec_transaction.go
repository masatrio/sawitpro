package postgres

import (
	"context"

	_ "github.com/lib/pq"
	"github.com/sawitpro/UserService/common"
)

func (c *Client) ExecTransaction(ctx context.Context, fn func(context.Context) error) error {
	tx, err := c.DB.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	ctx = context.WithValue(ctx, common.TX_KEY, tx)
	defer tx.Rollback()

	err = fn(ctx)
	if err != nil {
		return err
	}

	return tx.Commit()
}
