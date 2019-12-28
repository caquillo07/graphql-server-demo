package database

import (
	"github.com/jinzhu/gorm"
)

// Open creates a new connection with the given config
func Open(config Config) (*gorm.DB, error) {
	db, err := gorm.Open(config.Driver, config.ConnectionString)
	if err != nil {
		return nil, err
	}

	db.LogMode(config.Log)

	// Plural table names are lame
	db.SingularTable(true)

	// Do not allow update or delete to be called without a where clause.
	db.BlockGlobalUpdate(true)

	return db, nil
}