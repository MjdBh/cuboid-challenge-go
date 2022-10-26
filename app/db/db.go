package db

import (
	"fmt"

	"cuboid-challenge/app/config"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// CONN global gorm db connection.
var CONN *gorm.DB

// Connect set the gorm db connection.
func Connect() *gorm.DB {
	if CONN != nil {
		return CONN
	}

	var err error
	if CONN, err = gorm.Open(driver(), &gorm.Config{}); err != nil {
		panic(fmt.Errorf("failed to open the database connection. %w", err))
	}

	registerDBValidationsHooks(CONN)

	return CONN
}

func driver() gorm.Dialector {
	d := config.ENV.DBDriver
	if d == "sqlite" {
		return sqlite.Open(config.ENV.DBName)
	}

	panic(driverErr{d}.Error())
}

type driverErr struct{ name string }

func (d driverErr) Error() string {
	return fmt.Sprintf("DB driver '%s' not supported", d.name)
}
