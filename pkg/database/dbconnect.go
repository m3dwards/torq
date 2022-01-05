package database

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// PgConnect connects to a pg database
func PgConnect(user_name string, db_name string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", fmt.Sprintf("user=%s dbname=%s sslmode=disable", user_name, db_name))
	if err != nil {
		return nil, fmt.Errorf("internal/database/connect PgConnect: %v", err)
	}
	return db, nil
}
