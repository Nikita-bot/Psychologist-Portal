package dal

import (
	"log"
	"psyportal/entity"

	_ "modernc.org/sqlite"
)

func (r Repo) CreateSlot(slot entity.Slot) (int64, error) {
	result, err := r.db.Exec("INSERT INTO slots (time) VALUES (?)", slot.Time)
	if err != nil {
		log.Fatal(err.Error())
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
		return 0, nil
	}

	return id, nil
}

func (r Repo) CreateSlotRoom(slot entity.Slot) (int64, error) {
	result, err := r.db.Exec("INSERT INTO slots_room (time) VALUES (?)", slot.Time)
	if err != nil {
		log.Fatal(err.Error())
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
		return 0, nil
	}

	return id, nil
}

func (r Repo) GetAllSlots() ([]entity.Slot, error) {

	var s []entity.Slot

	rows, err := r.db.Query("SELECT id, time FROM slots ORDER BY time")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slot entity.Slot
		if err := rows.Scan(&slot.ID, &slot.Time); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		s = append(s, slot)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return s, nil
}

func (r Repo) GetAllSlotsRoom() ([]entity.Slot, error) {
	var s []entity.Slot

	rows, err := r.db.Query("SELECT id, time FROM slots_room ORDER BY time")
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var slot entity.Slot
		if err := rows.Scan(&slot.ID, &slot.Time); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		s = append(s, slot)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err.Error())
		return nil, err
	}

	return s, nil
}

func (r Repo) DeleteSlot(id int) error {
	_, err := r.db.Exec("DELETE FROM slots WHERE id = ?", id)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func (r Repo) DeleteSlotRoom(id int) error {
	_, err := r.db.Exec("DELETE FROM slots_room WHERE id = ?", id)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}

func (r Repo) GetSlotDays() ([]entity.SlotDate, error) {
	rows, err := r.db.Query("SELECT day, is_active FROM slots_day")
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	var days []entity.SlotDate
	for rows.Next() {
		var day entity.SlotDate
		var isActive int
		if err := rows.Scan(&day.Day, &isActive); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		day.IsActive = isActive == 1
		days = append(days, day)
	}

	return days, nil
}

func (r Repo) GetSlotDaysRoom() ([]entity.SlotDate, error) {
	rows, err := r.db.Query("SELECT day, is_active FROM slots_day_room")
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	var days []entity.SlotDate
	for rows.Next() {
		var day entity.SlotDate
		var isActive int
		if err := rows.Scan(&day.Day, &isActive); err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		day.IsActive = isActive == 1
		days = append(days, day)
	}

	return days, nil
}

func (r Repo) ChangeDays(days []entity.SlotDate) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	if _, err := tx.Exec("UPDATE slots_day SET is_active = 0"); err != nil {
		log.Fatal(err.Error())
		return err
	}

	for _, day := range days {
		if _, err := tx.Exec(
			"UPDATE slots_day SET is_active = ? WHERE day = ?",
			day.IsActive,
			day.Day,
		); err != nil {
			log.Fatal(err.Error())
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err.Error())
		tx.Rollback()
		return err
	}

	return nil
}

func (r Repo) ChangeDaysRoom(days []entity.SlotDate) error {
	tx, err := r.db.Begin()
	if err != nil {
		log.Fatal(err.Error())
		return err
	}

	if _, err := tx.Exec("UPDATE slots_day_room SET is_active = 0"); err != nil {
		log.Fatal(err.Error())
		return err
	}

	for _, day := range days {
		if _, err := tx.Exec(
			"UPDATE slots_day_room SET is_active = ? WHERE day = ?",
			day.IsActive,
			day.Day,
		); err != nil {
			log.Fatal(err.Error())
			return err
		}
	}

	if err := tx.Commit(); err != nil {
		log.Fatal(err.Error())
		tx.Rollback()
		return err
	}

	return nil
}
