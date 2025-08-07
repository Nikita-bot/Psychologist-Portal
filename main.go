package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"psyportal/config"
	"psyportal/dal"
	"psyportal/handler"
)

func main() {

	var err error
	db, err := sql.Open("sqlite", "./formdata.db")
	if err != nil {
		log.Fatal(err)
	}

	config.InitDB(db)

	repo := dal.InitRepo(db)

	handl := handler.InitHandler(repo)

	c := config.InitConfig()

	fs := http.FileServer(http.Dir("static"))

	mux := http.NewServeMux()

	mux.HandleFunc("GET /", handl.MainPage)

	mux.Handle("GET /static/", http.StripPrefix("/static/", fs))

	mux.HandleFunc("GET /room", handl.RoomPage)

	mux.HandleFunc("GET /event", handl.IndividualPage)

	mux.HandleFunc("POST /event", handl.PostEvent)

	mux.HandleFunc("GET /admin/event", adminAuth(handl.AdminPage, c))

	mux.HandleFunc("GET /admin/event/archive", adminAuth(handl.ArchivePage, c))

	mux.HandleFunc("GET /available-times", handl.GetTimes)

	mux.HandleFunc("GET /available-times-room", handl.GetTimesRoom)

	mux.HandleFunc("DELETE /admin/event", adminAuth(handl.DeleteEvent, c))

	mux.HandleFunc("GET /slots", adminAuth(handl.SlotPage, c))

	mux.HandleFunc("GET /slots_ind", handl.GetSlots)

	mux.HandleFunc("POST /slots", adminAuth(handl.PostSlot, c))

	mux.HandleFunc("DELETE /slots", adminAuth(handl.DeleteSlot, c))

	mux.HandleFunc("GET /slots_room", adminAuth(handl.GetSlotsRoom, c))

	mux.HandleFunc("POST /slots_room", adminAuth(handl.PostSlotRoom, c))

	mux.HandleFunc("DELETE /slots_room", adminAuth(handl.DeleteSlotRoom, c))

	mux.HandleFunc("GET /slots/days", handl.GetSlotDays)

	mux.HandleFunc("PATCH /slots/days", adminAuth(handl.PatchSlotDay, c))

	mux.HandleFunc("GET /slots_room/days", handl.GetSlotDaysRoom)

	mux.HandleFunc("PATCH /slots_room/days", adminAuth(handl.PatchSlotDayRoom, c))

	mux.HandleFunc("GET /employees", adminAuth(handl.EmployeePage, c))

	mux.HandleFunc("GET /employees_ind", handl.GetEmployees)

	mux.HandleFunc("POST /employees", adminAuth(handl.PostEmployee, c))

	mux.HandleFunc("DELETE /employees", adminAuth(handl.DeleteEmployee, c))

	mux.HandleFunc("PATCH /employees", adminAuth(handl.PatchEmployee, c))

	log.Print("Server started at :" + c.Port)
	port := fmt.Sprintf(":%s", c.Port)
	http.ListenAndServe(port, mux)
}

func adminAuth(next http.HandlerFunc, c *config.Config) http.HandlerFunc {
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
