package config

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	pkg "github.com/restuwahyu13/sso-and-cloudstorage/packages"
)

func Connection(driver string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Second*60))
	defer cancel()

	res, err := sqlx.ConnectContext(ctx, driver, pkg.GetString("PG_DSN"))
	return res, err
}
