package msi

import (
	"database/sql"
	"database/sql/driver"
	"log"

	"github.com/go-sql-driver/mysql"
)

func init() {
	sql.Register("mysqlMsi", NewMySQLWrapperDriver())
}

type MySQLWrapperDriver struct {
	delegate  *mysql.MySQLDriver
	msiClient *AccessTokenClient
}

func NewMySQLWrapperDriver() *MySQLWrapperDriver {
	return &MySQLWrapperDriver{
		&mysql.MySQLDriver{},
		NewAccessTokenClient("https://ossrdbms-aad.database.windows.net"),
	}
}

func (drv *MySQLWrapperDriver) Open(dsn string) (driver.Conn, error) {
	log.Printf("Opening connection: %v", dsn)

	atr, err := drv.msiClient.RequestToken()
	if err != nil {
		return nil, err
	}
	log.Printf("Got token: %v", atr)

	// Update the config with the MSI token for password
	config, err := mysql.ParseDSN(dsn)
	if err != nil {
		return nil, err
	}
	config.Passwd = atr.AccessToken
	config.AllowCleartextPasswords = true

	return drv.delegate.Open(config.FormatDSN())
}
