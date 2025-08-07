package dal

import (
	"log"
	"psyportal/entity"

	_ "modernc.org/sqlite"
)

func (r Repo) GetAllEmployees() ([]entity.Employees, error) {
	var emp []entity.Employees

	rows, err := r.db.Query(`
			SELECT id, name, tg 
			FROM employees 
		`)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e entity.Employees
		err := rows.Scan(&e.ID, &e.FullName, &e.Telegram)
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		emp = append(emp, e)
	}

	return emp, nil
}

func (r Repo) GetAllEmployeesName() ([]string, error) {
	var emp []string

	rows, err := r.db.Query(`
			SELECT name 
			FROM employees 
		`)
	if err != nil {
		log.Fatal(err.Error())
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var e string
		err := rows.Scan(&e)
		if err != nil {
			log.Fatal(err.Error())
			return nil, err
		}
		emp = append(emp, e)
	}

	return emp, nil
}

func (r Repo) CreateEmployee(emp entity.Employees) (int64, error) {
	result, err := r.db.Exec("INSERT INTO employees (name, tg) VALUES (?,?)", emp.FullName, emp.Telegram)
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		log.Fatal(err.Error())
		return 0, err
	}

	return id, nil
}

func (r Repo) ChangeEmployee(emp entity.Employees, id int) error {
	if _, err := r.db.Exec(
		"UPDATE employees SET name = ?, tg = ? WHERE id = ?",
		emp.FullName,
		emp.Telegram,
		id,
	); err != nil {
		log.Fatal(err.Error())
		return err
	}

	return nil
}

func (r Repo) DeleteEmployee(id int) error {
	_, err := r.db.Exec("DELETE FROM employees WHERE id = ?", id)
	if err != nil {
		log.Fatal(err.Error())
		return err
	}
	return nil
}
