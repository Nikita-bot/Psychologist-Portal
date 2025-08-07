package entity

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
