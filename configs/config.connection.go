package config

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	pkg "github.com/restuwahyu13/sso-and-cloudstorage/packages"
)

func Connection(driverName string) (*sqlx.DB, error) {
	ctx, cancel := context.WithTimeout(context.TODO(), time.Duration(time.Second*10))
	defer cancel()

	res, err := sqlx.ConnectContext(ctx, driverName, pkg.GetString("PG_DSN"))
	return res, err
}
