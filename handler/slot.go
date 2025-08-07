package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"psyportal/entity"
	"strconv"
	"text/template"
)

func (h Handler) SlotPage(w http.ResponseWriter, r *http.Request) {
	s, err := h.repo.GetAllSlots()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	tmpl, err := template.ParseFiles("template/slots.html")
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Template not found"}`, http.StatusInternalServerError)
	}

	err = tmpl.Execute(w, s)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to execute template"}`, http.StatusInternalServerError)
	}
}

func (h Handler) GetSlots(w http.ResponseWriter, r *http.Request) {

	s, err := h.repo.GetAllSlots()
	if err != nil {
		http.Error(w, `{"error": "Failed to get slots"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"slots":  s,
	})
}

func (h Handler) GetSlotsRoom(w http.ResponseWriter, r *http.Request) {
	slots, err := h.repo.GetAllSlotsRoom()
	if err != nil {
		http.Error(w, `{"error": "Failed to get room's slots"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"slots":  slots,
	})
}

func (h Handler) PostSlot(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}

	var slot entity.Slot
	err := json.NewDecoder(r.Body).Decode(&slot)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to parse data"}`, http.StatusBadRequest)
	}

	if slot.Time == "" {
		http.Error(w, `{"error":"Time is required"}`, http.StatusBadRequest)
	}

	id, err := h.repo.CreateSlot(slot)
	if err != nil {
		http.Error(w, `{"error": "Failed to create slot"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   id,
		"time": slot.Time,
	})
}

func (h Handler) PostSlotRoom(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error": "Method not allowed"}`, http.StatusMethodNotAllowed)
	}

	var slot entity.Slot
	err := json.NewDecoder(r.Body).Decode(&slot)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to parse data"}`, http.StatusBadRequest)
	}

	if slot.Time == "" {
		http.Error(w, `{"error":"Time is required"}`, http.StatusBadRequest)
	}

	id, err := h.repo.CreateSlotRoom(slot)
	if err != nil {
		http.Error(w, `{"error":"Failed to create room's slot"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":   id,
		"time": slot.Time,
	})
}

func (h Handler) DeleteSlot(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Missing event ID"}`, http.StatusBadRequest)
	}

	slotID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
	}

	err = h.repo.DeleteSlot(slotID)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete slot"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h Handler) DeleteSlotRoom(w http.ResponseWriter, r *http.Request) {
	idStr := r.URL.Query().Get("id")
	if idStr == "" {
		http.Error(w, `{"error": "Missing event ID"}`, http.StatusBadRequest)
	}

	slotID, err := strconv.Atoi(idStr)
	if err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Invalid event ID format"}`, http.StatusBadRequest)
	}

	err = h.repo.DeleteSlotRoom(slotID)
	if err != nil {
		http.Error(w, `{"error": "Failed to delete room's slot"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h Handler) GetSlotDays(w http.ResponseWriter, r *http.Request) {
	days, err := h.repo.GetSlotDays()
	if err != nil {
		http.Error(w, `{"error": "Failed to get day slots"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"days":   days,
	})
}

func (h Handler) GetSlotDaysRoom(w http.ResponseWriter, r *http.Request) {
	days, err := h.repo.GetSlotDaysRoom()
	if err != nil {
		http.Error(w, `{"error": "Failed to get room's day slots"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status": "success",
		"days":   days,
	})
}

func (h Handler) PatchSlotDay(w http.ResponseWriter, r *http.Request) {
	var days []entity.SlotDate

	if err := json.NewDecoder(r.Body).Decode(&days); err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to parse data"}`, http.StatusBadRequest)
	}

	log.Println("Days: ", days)

	err := h.repo.ChangeDays(days)
	if err != nil {
		http.Error(w, `{"error": "Failed to change day slots"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}

func (h Handler) PatchSlotDayRoom(w http.ResponseWriter, r *http.Request) {
	var days []entity.SlotDate

	if err := json.NewDecoder(r.Body).Decode(&days); err != nil {
		log.Fatal(err.Error())
		http.Error(w, `{"error": "Failed to parse data"}`, http.StatusBadRequest)
	}

	log.Println("Days: ", days)

	err := h.repo.ChangeDaysRoom(days)
	if err != nil {
		http.Error(w, `{"error": "Failed to change room's day slots"}`, http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "success"})
}
