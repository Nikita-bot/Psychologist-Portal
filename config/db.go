package config

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func InitDB(db *sql.DB) {

	sqlStmt := `
		CREATE TABLE IF NOT EXISTS consultations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			fio TEXT NOT NULL,
			phone TEXT NOT NULL,
			position TEXT NOT NULL,
			comment TEXT,
			meet_date TEXT NOT NULL,
			meet_time TEXT NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP
		);

		CREATE TABLE IF NOT EXISTS slots (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			time TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS slots_day (
			day TEXT UNIQUE,
			is_active INTEGER
		);

		CREATE TABLE IF NOT EXISTS slots_room (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			time TEXT NOT NULL
		);

		CREATE TABLE IF NOT EXISTS slots_day_room (
			day TEXT UNIQUE,
			is_active INTEGER
		);

		CREATE TABLE IF NOT EXISTS employees (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			name TEXT NOT NULL,
			tg TEXT NOT NULL
		);

		INSERT OR IGNORE INTO slots_day (day, is_active) VALUES
		('sunday', 0), 
		('tuesday', 0),
		('wednesday', 0),
		('thursday', 0),
		('friday', 0),
		('saturday', 0),
		('monday', 0);

		INSERT OR IGNORE INTO slots_day_room (day, is_active) VALUES
		('sunday', 0), 
		('tuesday', 0),
		('wednesday', 0),
		('thursday', 0),
		('friday', 0),
		('saturday', 0),
		('monday', 0);
    `

	_, err := db.Exec(sqlStmt)
	if err != nil {
		log.Fatalf("Ошибка создания таблицы: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('consultations') WHERE name='type'").Scan(&count)
	if err != nil {
		log.Printf("Ошибка проверки колонки: %v", err)
		return
	}

	if count == 0 {
		_, err = db.Exec("ALTER TABLE consultations ADD COLUMN type TEXT NOT NULL DEFAULT ''")
		if err != nil {
			log.Printf("Ошибка добавления колонки: %v", err)
		}
	}

	err = db.QueryRow("SELECT COUNT(*) FROM pragma_table_info('consultations') WHERE name='psyholog'").Scan(&count)
	if err != nil {
		log.Printf("Ошибка проверки колонки: %v", err)
		return
	}

	if count == 0 {
		_, err = db.Exec("ALTER TABLE consultations ADD COLUMN psyholog TEXT DEFAULT ''")
		if err != nil {
			log.Printf("Ошибка добавления колонки: %v", err)
		}
	}
}
