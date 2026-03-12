package cache

import (
	"database/sql"
)

type Cache struct {
	db *sql.DB
}
