package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"psyportal/entity"
	"strconv"
	"text/template"
)

func (h Handler) EmployeePage(w http.ResponseWriter, r *http.Request) {
	emp, err := h.repo.GetAllEmployees()
	if err != nil {
		http.Error(w, `{"error": "Failed to get employee"}`, http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles("template/users.html")
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Template not found"}`, http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, emp)
	if err != nil {
		log.Println(err.Error())
		http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
	}

}

func (h Handler) GetEmployees(w http.ResponseWriter, r *http.Request) {
	empname, err := h.repo.GetAllEmployeesName()
	if err != nil {
		http.Error(w, `{"error": "Failed to get employee"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(empname)
}

func (h Handler) PostEmployee(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
	}

	var emp entity.Employees
	err := json.NewDecoder(r.Body).Decode(&emp)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to parse data"}`, http.StatusBadRequest)
	}

	id, err := h.repo.CreateEmployee(emp)
	if err != nil {
		http.Error(w, `{"error": "Failed to create employee"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id": id,
	})
}

func (h Handler) DeleteEmployee(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Missing employee ID"}`, http.StatusBadRequest)
	}

	empId, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Invalid employee ID format"}`, http.StatusBadRequest)
	}

	err = h.repo.DeleteEmployee(empId)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete employee"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h Handler) PatchEmployee(w http.ResponseWriter, r *http.Request) {
	var emp entity.Employees

	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error": "Missing employee ID"}`, http.StatusBadRequest)
	}

	empID, err := strconv.Atoi(id)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Invalid employee ID format"}`, http.StatusBadRequest)
	}

	if err := json.NewDecoder(r.Body).Decode(&emp); err != nil {
		log.Println(err.Error())
		http.Error(w, `{"error": "Failed to parse data"}`, http.StatusBadRequest)
	}

	log.Println("Emp: ", emp)

	err = h.repo.ChangeEmployee(emp, empID)
	if err != nil {
		http.Error(w, `{"error": "Failed to change employee"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
