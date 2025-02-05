package database

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/glebarez/sqlite"
	_ "gorm.io/gorm"
	"log"
	"os"
)

type Database struct {
	Name string
	db   *sql.DB
}

func New(name string) (*Database, error) {
	_, err := os.ReadDir("database")
	if err != nil {
		if os.IsNotExist(err) {
			err = os.Mkdir("database", 0755)
			if err != nil {
				return nil, err
			}
		} else {
			return nil, err
		}
	}

	db, err := sql.Open("sqlite", fmt.Sprintf("database/%s", name))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	err = createTableIfNotExists(db)
	if err != nil {
		return nil, err
	}

	// Check database integrity
	err = checkDatabaseIntegrity(db)
	if err != nil {
		return nil, err
	}

	return &Database{Name: name, db: db}, nil
}

func (d *Database) Set(key string, value interface{}) error {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = d.db.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, jsonValue)
	return err
}

func (d *Database) Get(key string) (interface{}, error) {
	var jsonValue string
	err := d.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&jsonValue)
	if err != nil {
		return nil, err
	}

	var result interface{}
	err = json.Unmarshal([]byte(jsonValue), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (d *Database) Delete(key string) error {
	_, err := d.db.Exec("DELETE FROM settings WHERE key = ?", key)
	return err
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Has(key string) (bool, error) {
	var value string
	err := d.db.QueryRow("SELECT value FROM settings WHERE key = ?", key).Scan(&value)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (d *Database) GetAll() (map[string]interface{}, error) {
	rows, err := d.db.Query("SELECT key, value FROM settings")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	result := make(map[string]interface{})
	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			return nil, err
		}
		var resultValue interface{}
		err := json.Unmarshal([]byte(value), &resultValue)
		if err != nil {
			return nil, err
		}
		result[key] = resultValue
	}

	return result, nil
}

func createTableIfNotExists(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS settings (
		key TEXT PRIMARY KEY,
		value BLOB
	);`
	_, err := db.Exec(query)
	return err
}

func (d *Database) List() ([]string, error) {
	var keys []string
	rows, err := d.db.Query("SELECT key FROM settings")
	if err != nil {
		return nil, err
	}
	defer func(rows *sql.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("Error closing rows: %v", err)
		}
	}(rows)

	for rows.Next() {
		var key string
		if err := rows.Scan(&key); err != nil {
			return nil, err
		}
		keys = append(keys, key)
	}

	return keys, nil
}

func checkDatabaseIntegrity(db *sql.DB) error {
	var integrityCheck string
	err := db.QueryRow("PRAGMA integrity_check").Scan(&integrityCheck)
	if err != nil {
		return err
	}
	if integrityCheck != "ok" {
		return fmt.Errorf("database integrity check failed: %s", integrityCheck)
	}
	return nil
}
