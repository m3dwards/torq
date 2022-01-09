package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PgConnect connects to a pg database
func PgConnect(dbName, userName, dbPassword string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres",
		fmt.Sprintf("user=%s dbname=%s password=%s sslmode=disable", userName, dbName, dbPassword))
	if err != nil {
		return nil, fmt.Errorf("internal/database/connect PgConnect: %v", err)
	}
	return db, nil
}
