package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"psyportal/config"
	"psyportal/entity"
	"strconv"
	"text/template"
)

func (h Handler) MainPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./template/main.html")
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Template not found"}`, http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
	}
}

func (h Handler) RoomPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./template/room.html")
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Template not found"}`, http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
	}
}

func (h Handler) IndividualPage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./template/individual.html")
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Template not found"}`, http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
	}
}

func (h Handler) AdminPage(w http.ResponseWriter, r *http.Request) {
	days, err := h.repo.GetFormsByDays()
	if err != nil {
		http.Error(w, `{"error": "Failed to get cards"}`, http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles("template/admin.html")
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Template not found"}`, http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, days)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
	}
}

func (h Handler) ArchivePage(w http.ResponseWriter, r *http.Request) {
	archiveData, err := h.repo.GetArchiveConsultations()
	if err != nil {
		http.Error(w, `{"error": "Failed to get cards"}`, http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles("template/archive.html")
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Template not found"}`, http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, archiveData)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
	}
}

func (h Handler) PostEvent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}

	var form entity.UserForm
	err := json.NewDecoder(r.Body).Decode(&form)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to parse data"}`, http.StatusBadRequest)
	}

	h.repo.SaveFormToDB(form)

	w.Header().Set("Content-Type", "application/json")

	response := map[string]interface{}{
		"status":  "success",
		"message": "Форма успешно отправлена",
		"data":    form,
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to send data"}`, http.StatusInternalServerError)
	}
}

func (h Handler) DeleteEvent(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Missing event ID"}`, http.StatusBadRequest)
	}

	eventID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
	}

	if err := h.repo.DeleteData(eventID); err != nil {
		http.Error(w, `{"error": "Failed to delete card"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h Handler) GetTimes(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, `{"error": "Date parameter is required"}`, http.StatusBadRequest)
	}

	bookedTimes, err := h.repo.GetAvailableTimes(date, config.INDIVIDUAL)
	if err != nil {
		http.Error(w, `{"error": "Failed to get available times"}`, http.StatusInternalServerError)
	}

	allSlots, _ := h.repo.GetAllSlots()

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
}

func (h Handler) GetTimesRoom(w http.ResponseWriter, r *http.Request) {
	date := r.URL.Query().Get("date")
	if date == "" {
		http.Error(w, `{"error": "Date parameter is required"}`, http.StatusBadRequest)
	}

	bookedTimes, err := h.repo.GetAvailableTimesRoom(date, config.ROOM)
	if err != nil {
		http.Error(w, `{"error": "Failed to get available times for room"}`, http.StatusInternalServerError)
	}

	allSlots, _ := h.repo.GetAllSlotsRoom()

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
}
