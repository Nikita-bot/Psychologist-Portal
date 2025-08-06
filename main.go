package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"text/template"
	"time"

	"github.com/goloop/env"
	_ "modernc.org/sqlite"
)

type UserForm struct {
	ID        int    `json:"id" db:"id"`
	FIO       string `json:"fio" db:"fio"`
	Phone     string `json:"phone" db:"phone"`
	Position  string `json:"position" db:"position"`
	Date      string `json:"meet_date" db:"meet_date"`
	TypeEvent string `json:"type" db:"type"`
	Time      string `json:"meet_time" db:"meet_time"`
	Psyholog  string `json:"psyholog" db:"psyholog"`
	Comment   string `json:"comment" db:"comment"`
}

type DayView struct {
	Date      string
	CountInd  int
	CountRoom int
	Forms     []UserForm
}

type Config struct {
	Login    string `env:"LOGIN"`
	Password string `env:"PASS"`
	Port     string `env:"PORT"`
}

type Slot struct {
	ID   int    `json:"id" db:"id"`
	Time string `json:"time" db:"time"`
}

type SlotDate struct {
	Day      string `json:"day" db:"day"`
	IsActive bool   `json:"is_active" db:"is_active"`
}

type Employees struct {
	ID       int    `json:"id" db:"id"`
	FullName string `json:"name" db:"name"`
	Telegram string `json:"tg" db:"tg"`
}

const (
	INDIVIDUAL = "Индивидуальная консультация"
	ROOM       = "Комната психологической разгрузки"
)

func initDB(db *sql.DB) {

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

func initConfig() *Config {
	if err := env.Load(".env"); err != nil {
		log.Fatal(err)
	}

	var cfg Config
	if err := env.Unmarshal("", &cfg); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("%+v\n", cfg)

	return &cfg
}

func saveFormToDB(db *sql.DB, form UserForm) error {
	log.Printf("Данные формы:\nFIO: %s\nPhone: %s\nPosition: %s\nPsyholog: %s\nComment: %s\nType: %s\nDate: %s\nTime: %s\n",
		form.FIO, form.Phone, form.Position, form.Psyholog,
		form.Comment, form.TypeEvent, form.Date, form.Time)
	stmt, err := db.Prepare(`
		INSERT INTO consultations(fio, phone, position, psyholog, comment, type, meet_date, meet_time)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Println(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(form.FIO, form.Phone, form.Position, form.Psyholog, form.Comment, form.TypeEvent, form.Date, form.Time)
	if err != nil {
		log.Println(err.Error())
	}
	return err
}

func getCountByType(db *sql.DB, date string) []int {
	var individualCount, roomCount int

	// Запрос для индивидуальных консультаций
	err := db.QueryRow(`
        SELECT COUNT(*) 
        FROM consultations 
        WHERE meet_date = ? AND type = ?
    `, date, INDIVIDUAL).Scan(&individualCount)
	if err != nil {
		return nil
	}

	// Запрос для комнаты разгрузки
	err = db.QueryRow(`
        SELECT COUNT(*) 
        FROM consultations 
        WHERE meet_date = ? AND type = ?
    `, date, ROOM).Scan(&roomCount)
	if err != nil {
		return nil
	}

	return []int{individualCount, roomCount}
}

func getFormsByDays(db *sql.DB) ([]DayView, error) {
	// Получаем текущую дату и 7 дней вперед
	now := time.Now()
	var days []DayView

	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i).Format("2006-01-02")
		days = append(days, DayView{Date: date})
	}

	// Получаем все записи
	rows, err := db.Query(`
        SELECT id, fio, phone, position, psyholog, comment, type, meet_date, meet_time 
        FROM consultations 
        WHERE date(meet_date) BETWEEN date('now') AND date('now', '+7 days')
        ORDER BY meet_date, meet_time
    `)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Группируем записи по датам
	dateMap := make(map[string][]UserForm)
	for rows.Next() {
		var f UserForm
		err := rows.Scan(&f.ID, &f.FIO, &f.Phone, &f.Position, &f.Psyholog, &f.Comment, &f.TypeEvent, &f.Date, &f.Time)
		if err != nil {
			return nil, err
		}
		dateMap[f.Date] = append(dateMap[f.Date], f)
	}

	var counts []int

	// Сопоставляем с нашими днями
	for i, day := range days {
		if forms, ok := dateMap[day.Date]; ok {
			days[i].Forms = forms
			counts = getCountByType(db, day.Date)
			days[i].CountInd = counts[0]
			days[i].CountRoom = counts[1]
		}
	}

	return days, nil
}

func getArchiveConsultations(db *sql.DB) ([]DayView, error) {

	query := `
        SELECT 
            DATE(meet_date) as date,
			id,
            fio, 
            phone, 
            position,
			psyholog,
			type, 
            meet_date, 
            meet_time, 
            comment
        FROM consultations
        WHERE meet_date < CURRENT_DATE
        ORDER BY meet_date DESC, meet_time DESC
		LIMIT 200
    `

	rows, err := db.Query(query)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса архивных данных: %v", err)
	}
	defer rows.Close()

	var archive []DayView
	var currentDate string
	var currentDay *DayView

	for rows.Next() {
		var (
			dateStr string
			form    UserForm
		)

		if err := rows.Scan(
			&dateStr,
			&form.ID,
			&form.FIO,
			&form.Phone,
			&form.Position,
			&form.Psyholog,
			&form.TypeEvent,
			&form.Date,
			&form.Time,
			&form.Comment,
		); err != nil {
			return nil, fmt.Errorf("ошибка сканирования данных: %v", err)
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			return nil, fmt.Errorf("ошибка парсинга даты: %v", err)
		}
		formattedDate := date.Format("02.01.2006")

		if currentDate != dateStr {
			currentDate = dateStr
			archive = append(archive, DayView{
				Date:  formattedDate,
				Forms: []UserForm{},
			})
			currentDay = &archive[len(archive)-1]
		}

		currentDay.Forms = append(currentDay.Forms, form)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при обработке результатов: %v", err)
	}

	return archive, nil
}

func deleteData(db *sql.DB, eventID int) error {

	_, err := db.Exec("DELETE FROM consultations WHERE id = ?", eventID)
	if err != nil {
		return fmt.Errorf("failed to delete event: %v", err)
	}

	return nil
}

func adminAuth(next http.HandlerFunc, c *Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Получаем логин и пароль из Basic Auth
		username, password, ok := r.BasicAuth()

		if !ok || username != c.Login || password != c.Password {
			// Запрашиваем авторизацию
			w.Header().Set("WWW-Authenticate", `Basic realm="Restricted"`)
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		next(w, r)
	}
}

func getAllSlots(db *sql.DB) ([]Slot, error) {

	var s []Slot

	rows, err := db.Query("SELECT id, time FROM slots ORDER BY time")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slot Slot
		if err := rows.Scan(&slot.ID, &slot.Time); err != nil {
			return nil, fmt.Errorf("failed to scan slot: %w", err)
		}
		s = append(s, slot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return s, nil
}

func getAllSlotsRoom(db *sql.DB) ([]Slot, error) {

	var s []Slot

	rows, err := db.Query("SELECT id, time FROM slots_room ORDER BY time")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slot Slot
		if err := rows.Scan(&slot.ID, &slot.Time); err != nil {
			return nil, fmt.Errorf("failed to scan slot: %w", err)
		}
		s = append(s, slot)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows error: %w", err)
	}

	return s, nil
}

func main() {

	var err error
	db, err := sql.Open("sqlite", "./formdata.db")
	if err != nil {
		log.Fatal(err)
	}

	initDB(db)

	c := initConfig()

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./template/main.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("GET /room", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./template/room.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("GET /event", func(w http.ResponseWriter, r *http.Request) {
		tmpl, err := template.ParseFiles("./template/individual.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})

	mux.HandleFunc("POST /event", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			log.Println(err.Error())
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var form UserForm
		err := json.NewDecoder(r.Body).Decode(&form)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		saveFormToDB(db, form)

		w.Header().Set("Content-Type", "application/json")

		response := map[string]interface{}{
			"status":  "success",
			"message": "Форма успешно отправлена",
			"data":    form,
		}

		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	})

	mux.HandleFunc("GET /admin/event", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		days, err := getFormsByDays(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("template/admin.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, days)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}, c))

	mux.HandleFunc("GET /admin/event/archive", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		// Получаем архивные данные (пример функции)
		archiveData, err := getArchiveConsultations(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("template/archive.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, archiveData)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}, c))

	mux.HandleFunc("GET /available-times", func(w http.ResponseWriter, r *http.Request) {

		date := r.URL.Query().Get("date")
		if date == "" {
			http.Error(w, "Date parameter is required", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
        SELECT meet_time 
        FROM consultations 
        WHERE meet_date = ? 
		AND type = ?
        ORDER BY meet_time`, date, INDIVIDUAL)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bookedTimes []string
		for rows.Next() {
			var time string
			if err := rows.Scan(&time); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			bookedTimes = append(bookedTimes, time)
		}

		allSlots, _ := getAllSlots(db)

		availableSlots := make([]string, 0)
		for _, slot := range allSlots {
			isBooked := false
			for _, booked := range bookedTimes {
				if booked == slot.Time {
					isBooked = true
					break
				}
			}
			if !isBooked {
				availableSlots = append(availableSlots, slot.Time)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(availableSlots)
	})

	mux.HandleFunc("GET /available-times-room", func(w http.ResponseWriter, r *http.Request) {

		date := r.URL.Query().Get("date")
		if date == "" {
			http.Error(w, "Date parameter is required", http.StatusBadRequest)
			return
		}

		rows, err := db.Query(`
        SELECT COUNT(*), meet_time 
        FROM consultations 
        WHERE meet_date = ?
		AND type = ?
		GROUP BY meet_time
        ORDER BY meet_time`, date, ROOM)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		var bookedTimes = make(map[string]int)
		for rows.Next() {
			var time string
			var count int
			if err := rows.Scan(&count, &time); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			bookedTimes[time] = count
		}

		allSlots, _ := getAllSlotsRoom(db)

		availableSlots := make([]string, 0)
		for _, slot := range allSlots {
			isBooked := false
			for time, count := range bookedTimes {
				if time == slot.Time && count >= 7 {
					isBooked = true
					break
				}
			}
			if !isBooked {
				availableSlots = append(availableSlots, slot.Time)
			}
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(availableSlots)
	})

	mux.HandleFunc("DELETE /admin/event", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error": "Missing event ID"}`, http.StatusBadRequest)
			return
		}

		eventID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
			return
		}

		if err := deleteData(db, eventID); err != nil {
			http.Error(w, fmt.Sprintf(`{"error": "%v"}`, err), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}, c))

	mux.HandleFunc("GET /slots", adminAuth(func(w http.ResponseWriter, r *http.Request) {

		var s []Slot

		s, err := getAllSlots(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		tmpl, err := template.ParseFiles("template/slots.html")
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, s)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}, c))

	mux.HandleFunc("GET /slots_ind", func(w http.ResponseWriter, r *http.Request) {

		var s []Slot

		s, err := getAllSlots(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"slots":  s,
		})
	})

	mux.HandleFunc("POST /slots", adminAuth(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var slot Slot
		err := json.NewDecoder(r.Body).Decode(&slot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if slot.Time == "" {
			http.Error(w, `{"error":"Time is required"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO slots (time) VALUES (?)", slot.Time)
		if err != nil {
			http.Error(w, `{"error":"Failed to create slot"}`, http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Failed to get slot ID"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   id,
			"time": slot.Time,
		})

	}, c))

	mux.HandleFunc("DELETE /slots", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error": "Missing event ID"}`, http.StatusBadRequest)
			return
		}

		slotID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM slots WHERE id = ?", slotID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error"})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}, c))

	mux.HandleFunc("GET /slots_room", adminAuth(func(w http.ResponseWriter, r *http.Request) {

		slots, err := getAllSlotsRoom(db)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"slots":  slots,
		})
	}, c))

	mux.HandleFunc("POST /slots_room", adminAuth(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var slot Slot
		err := json.NewDecoder(r.Body).Decode(&slot)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if slot.Time == "" {
			http.Error(w, `{"error":"Time is required"}`, http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO slots_room (time) VALUES (?)", slot.Time)
		if err != nil {
			http.Error(w, `{"error":"Failed to create slot"}`, http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Failed to get slot ID"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":   id,
			"time": slot.Time,
		})

	}, c))

	mux.HandleFunc("DELETE /slots_room", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error": "Missing event ID"}`, http.StatusBadRequest)
			return
		}

		slotID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM slots_room WHERE id = ?", slotID)
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(map[string]string{"status": "error"})
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}, c))

	mux.HandleFunc("GET /slots/days", (func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT day, is_active FROM slots_day")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			return
		}
		defer rows.Close()

		var days []SlotDate
		for rows.Next() {
			var day SlotDate
			var isActive int
			if err := rows.Scan(&day.Day, &isActive); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to scan day"})
				return
			}
			day.IsActive = isActive == 1
			days = append(days, day)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"days":   days,
		})
	}))

	mux.HandleFunc("PATCH /slots/days", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		// Ожидаем массив объектов {day, active}
		var days []SlotDate

		if err := json.NewDecoder(r.Body).Decode(&days); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		log.Println("Days: ", days)

		tx, err := db.Begin()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			return
		}

		if _, err := tx.Exec("UPDATE slots_day SET is_active = 0"); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to reset days"})
			return
		}

		for _, day := range days {
			if _, err := tx.Exec(
				"UPDATE slots_day SET is_active = ? WHERE day = ?",
				day.IsActive,
				day.Day,
			); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Failed to update day: " + day.Day,
				})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{"error": "Transaction failed"})
			tx.Rollback()
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}, c))

	mux.HandleFunc("GET /slots_room/days", (func(w http.ResponseWriter, r *http.Request) {
		rows, err := db.Query("SELECT day, is_active FROM slots_day_room")
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			return
		}
		defer rows.Close()

		var days []SlotDate
		for rows.Next() {
			var day SlotDate
			var isActive int
			if err := rows.Scan(&day.Day, &isActive); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				json.NewEncoder(w).Encode(map[string]string{"error": "Failed to scan day"})
				return
			}
			day.IsActive = isActive == 1
			days = append(days, day)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]interface{}{
			"status": "success",
			"days":   days,
		})
	}))

	mux.HandleFunc("PATCH /slots_room/days", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		// Ожидаем массив объектов {day, active}
		var days []SlotDate

		if err := json.NewDecoder(r.Body).Decode(&days); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{
				"error":   "Invalid request body",
				"details": err.Error(),
			})
			return
		}

		log.Println("Days ROOM: ", days)

		tx, err := db.Begin()
		if err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{"error": "Database error"})
			return
		}

		if _, err := tx.Exec("UPDATE slots_day_room SET is_active = 0"); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{"error": "Failed to reset days"})
			return
		}

		for _, day := range days {
			if _, err := tx.Exec(
				"UPDATE slots_day_room SET is_active = ? WHERE day = ?",
				day.IsActive,
				day.Day,
			); err != nil {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				log.Println(err.Error())
				json.NewEncoder(w).Encode(map[string]string{
					"error": "Failed to update day: " + day.Day,
				})
				return
			}
		}

		if err := tx.Commit(); err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusInternalServerError)
			log.Println(err.Error())
			json.NewEncoder(w).Encode(map[string]string{"error": "Transaction failed"})
			tx.Rollback()
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}, c))

	mux.HandleFunc("GET /employees", adminAuth(func(w http.ResponseWriter, r *http.Request) {

		var emp []Employees

		rows, err := db.Query(`
			SELECT id, name, tg 
			FROM employees 
		`)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var e Employees
			err := rows.Scan(&e.ID, &e.FullName, &e.Telegram)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			emp = append(emp, e)
		}

		tmpl, err := template.ParseFiles("template/users.html")
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, emp)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}

	}, c))

	mux.HandleFunc("GET /employees_ind", func(w http.ResponseWriter, r *http.Request) {

		var empname []string

		rows, err := db.Query(`
			SELECT name 
			FROM employees 
		`)
		if err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer rows.Close()

		for rows.Next() {
			var e string
			err := rows.Scan(&e)
			if err != nil {
				log.Println(err.Error())
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			empname = append(empname, e)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(empname)
	})

	mux.HandleFunc("POST /employees", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
			return
		}

		var emp Employees
		err := json.NewDecoder(r.Body).Decode(&emp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		result, err := db.Exec("INSERT INTO employees (name, tg) VALUES (?,?)", emp.FullName, emp.Telegram)
		if err != nil {
			http.Error(w, `{"error":"Failed to create slot"}`, http.StatusInternalServerError)
			return
		}

		id, err := result.LastInsertId()
		if err != nil {
			http.Error(w, `{"error":"Failed to get slot ID"}`, http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": id,
		})
	}, c))

	mux.HandleFunc("DELETE /employees", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		idStr := r.URL.Query().Get("id")
		if idStr == "" {
			http.Error(w, `{"error": "Missing event ID"}`, http.StatusBadRequest)
			return
		}

		eventID, err := strconv.Atoi(idStr)
		if err != nil {
			http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
			return
		}

		_, err = db.Exec("DELETE FROM employees WHERE id = ?", eventID)
		if err != nil {
			http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}, c))

	mux.HandleFunc("PATCH /employees", adminAuth(func(w http.ResponseWriter, r *http.Request) {
		var emp Employees

		id := r.URL.Query().Get("id")
		if id == "" {
			log.Println(err.Error())
			http.Error(w, "ID сотрудника не указан", http.StatusBadRequest)
			return
		}

		if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		log.Println("Emp: ", emp)

		if _, err := db.Exec(
			"UPDATE employees SET name = ?, tg = ? WHERE id = ?",
			emp.FullName,
			emp.Telegram,
			id,
		); err != nil {
			log.Println(err.Error())
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(map[string]string{"status": "success"})
	}, c))

	log.Print("Server started at :" + c.Port)
	port := fmt.Sprintf(":%s", c.Port)
	http.ListenAndServe(port, mux)
}
