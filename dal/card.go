package dal

import (
	"fmt"
	"log"
	"psyportal/config"
	"psyportal/entity"
	"time"

	_ "modernc.org/sqlite"
)

func (r Repo) SaveFormToDB(form entity.UserForm) error {
	log.Printf("Данные формы:\nFIO: %s\nPhone: %s\nPosition: %s\nPsyholog: %s\nComment: %s\nType: %s\nDate: %s\nTime: %s\n",
		form.FIO, form.Phone, form.Position, form.Psyholog,
		form.Comment, form.TypeEvent, form.Date, form.Time)
	stmt, err := r.db.Prepare(`
		INSERT INTO consultations(fio, phone, position, psyholog, comment, type, meet_date, meet_time)
		VALUES(?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(form.FIO, form.Phone, form.Position, form.Psyholog, form.Comment, form.TypeEvent, form.Date, form.Time)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func (r Repo) GetCountByType(date string) []int {
	var individualCount, roomCount int

	// Запрос для индивидуальных консультаций
	err := r.db.QueryRow(`
        SELECT COUNT(*) 
        FROM consultations 
        WHERE meet_date = ? AND type = ?
    `, date, config.INDIVIDUAL).Scan(&individualCount)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	// Запрос для комнаты разгрузки
	err = r.db.QueryRow(`
        SELECT COUNT(*) 
        FROM consultations 
        WHERE meet_date = ? AND type = ?
    `, date, config.ROOM).Scan(&roomCount)
	if err != nil {
		log.Fatal(err.Error())
		return nil
	}

	return []int{individualCount, roomCount}
}

func (r Repo) GetFormsByDays() ([]entity.DayView, error) {
	// Получаем текущую дату и 7 дней вперед
	now := time.Now()
	var days []entity.DayView

	for i := 0; i < 7; i++ {
		date := now.AddDate(0, 0, i).Format("2006-01-02")
		days = append(days, entity.DayView{Date: date})
	}

	// Получаем все записи
	rows, err := r.db.Query(`
        SELECT id, fio, phone, position, psyholog, comment, type, meet_date, meet_time 
        FROM consultations 
        WHERE date(meet_date) BETWEEN date('now') AND date('now', '+7 days')
        ORDER BY meet_date, meet_time
    `)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	// Группируем записи по датам
	dateMap := make(map[string][]entity.UserForm)
	for rows.Next() {
		var f entity.UserForm
		err := rows.Scan(&f.ID, &f.FIO, &f.Phone, &f.Position, &f.Psyholog, &f.Comment, &f.TypeEvent, &f.Date, &f.Time)
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		dateMap[f.Date] = append(dateMap[f.Date], f)
	}

	var counts []int

	// Сопоставляем с нашими днями
	for i, day := range days {
		if forms, ok := dateMap[day.Date]; ok {
			days[i].Forms = forms
			counts = r.GetCountByType(day.Date)
			days[i].CountInd = counts[0]
			days[i].CountRoom = counts[1]
		}
	}

	return days, nil
}

func (r Repo) GetArchiveConsultations() ([]entity.DayView, error) {

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

	rows, err := r.db.Query(query)
	if err != nil {
		log.Fatal(err.Error())
		return nil, fmt.Errorf("ошибка запроса архивных данных: %v", err)
	}
	defer rows.Close()

	var archive []entity.DayView
	var currentDate string
	var currentDay *entity.DayView

	for rows.Next() {
		var (
			dateStr string
			form    entity.UserForm
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
			log.Fatal(err.Error())
			return nil, err
		}

		date, err := time.Parse("2006-01-02", dateStr)
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		formattedDate := date.Format("02.01.2006")

		if currentDate != dateStr {
			currentDate = dateStr
			archive = append(archive, entity.DayView{
				Date:  formattedDate,
				Forms: []entity.UserForm{},
			})
			currentDay = &archive[len(archive)-1]
		}

		currentDay.Forms = append(currentDay.Forms, form)
	}

	if err := rows.Err(); err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return archive, nil
}

func (r Repo) DeleteData(eventID int) error {

	_, err := r.db.Exec("DELETE FROM consultations WHERE id = ?", eventID)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (r Repo) GetAvailableTimes(date string, typeEvent string) ([]string, error) {

	rows, err := r.db.Query(`
        SELECT meet_time 
        FROM consultations 
        WHERE meet_date = ? 
		AND type = ?
        ORDER BY meet_time`, date, typeEvent)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	var bookedTimes []string
	for rows.Next() {
		var time string
		if err := rows.Scan(&time); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		bookedTimes = append(bookedTimes, time)
	}

	return bookedTimes, nil

}

func (r Repo) GetAvailableTimesRoom(date string, typeEvent string) (map[string]int, error) {

	rows, err := r.db.Query(`
        SELECT COUNT(*), meet_time 
        FROM consultations 
        WHERE meet_date = ?
		AND type = ?
		GROUP BY meet_time
        ORDER BY meet_time`, date, config.ROOM)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	defer rows.Close()

	var bookedTimes = make(map[string]int)
	for rows.Next() {
		var time string
		var count int
		if err := rows.Scan(&count, &time); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		bookedTimes[time] = count
	}

	return bookedTimes, nil

}
