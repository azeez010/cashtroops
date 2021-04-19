package database

import (
	"github.com/adigunhammedolalekan/cashtroops/types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

func Open(connectUri string) (*gorm.DB, error) {
	database, err := gorm.Open("postgres", connectUri)
	if err != nil {
		return nil, err
	}
	runMigration(database)
	return database, nil
}

func runMigration(db *gorm.DB) {
	db.Debug().AutoMigrate(&types.User{},
		&types.Verification{},
		&types.PasswordResetToken{},
		&types.Beneficiary{},
		&types.Balance{},
		&types.Address{},
		&types.Hook{},
		&types.Rate{},
		&types.Payment{})
}
