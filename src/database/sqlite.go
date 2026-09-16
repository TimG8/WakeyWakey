package database

import (
	"database/sql"
	"fmt"
	"sync"
	"wakeywakey/utils"

	_ "modernc.org/sqlite"
)

var (
	db   *sql.DB
	once sync.Once
)

// Init opens the SQLite database once with sensible defaults.
func Init(dbPath string) error {
	var err error

	once.Do(func() {
		dsn := fmt.Sprintf("file:%s?_pragma=busy_timeout(5000)&_pragma=foreign_keys(1)", dbPath)

		db, err = sql.Open("sqlite", dsn)
		if err != nil {
			return
		}

		err = db.Ping()
	})

	db.Exec("CREATE TABLE IF NOT EXISTS wakes (id INTEGER PRIMARY KEY AUTOINCREMENT, user_id TEXT NOT NULL, global INTEGER NOT NULL, alias TEXT NOT NULL, ip_address TEXT NOT NULL, mac_address TEXT NOT NULL)")

	return err
}

func AddWakeEntry(userId string, global bool, alias string, ip string, mac string) error {
	if db == nil {
		return fmt.Errorf("Database not initialized")
	}

	if (!utils.IsValidMacAddress(mac)) {
		return fmt.Errorf("Invalid MAC address format")
	}

	// check if alias already exists for user
	var exists int
	var err error
	err = db.QueryRow("SELECT COUNT(*) FROM wakes WHERE user_id = ? AND alias = ?", userId, alias).Scan(&exists)
	if err != nil || exists > 0 {
		return fmt.Errorf("%s", "Alias " + alias + " already exists")
	}

	_, err = db.Exec("INSERT INTO wakes (user_id, global, alias, ip_address, mac_address) VALUES (?, ?, ?, ?, ?)", userId, global, alias, ip, mac)
	return err
}

func RemoveWakeEntryByAlias(userId string, alias string) error {
	if db == nil {
		return fmt.Errorf("Database not initialized")
	}

	_, err := db.Exec("DELETE FROM wakes WHERE user_id = ? AND alias = ?", userId, alias)
	return err
}

type WakeEntry struct {
	Alias      string
	IpAddress  string
	MacAddress string
}

func GetAllEntriesByUserId(userId string) ([]WakeEntry, error) {
	if db == nil {
		return nil, fmt.Errorf("Database not initialized")
	}

	rows, err := db.Query("SELECT alias, ip_address, mac_address FROM wakes WHERE user_id = ? OR global = 1", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var entries []WakeEntry
	for rows.Next() {
		var entry WakeEntry
		if err := rows.Scan(&entry.Alias, &entry.IpAddress, &entry.MacAddress); err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}

	return entries, nil
}

func GetMacByAlias(userId string, alias string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("Database not initialized")
	}

	var mac string
	err := db.QueryRow("SELECT mac_address FROM wakes WHERE (user_id = ? OR global = 1) AND alias = ?", userId, alias).Scan(&mac)
	if err != nil {
		return "", err
	}

	return mac, nil
}

func GetIpByAlias(userId string, alias string) (string, error) {
	if db == nil {
		return "", fmt.Errorf("Database not initialized")
	}

	var ip string
	err := db.QueryRow("SELECT ip_address FROM wakes WHERE (user_id = ? OR global = 1) AND alias = ?", userId, alias).Scan(&ip)
	if err != nil {
		return "", err
	}

	return ip, nil
}