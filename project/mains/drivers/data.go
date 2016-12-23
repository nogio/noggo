package drivers

import (
    "github.com/nogio/noggo"
    "github.com/nogio/noggo/driver/data-postgres"
    "github.com/nogio/noggo/driver/data-mysql"
)

func init() {
    noggo.Driver("postgres", data_postgres.Driver())
    noggo.Driver("mysql", data_mysql.Driver())
}
