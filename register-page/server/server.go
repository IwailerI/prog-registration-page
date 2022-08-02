package server

import (
	"log"
	"net/http"
	"os"
	"time"
	"web-server/db_sqlite"
	"web-server/registrationform"
)

type Database interface {
	Open() error
	Close()
	Create() error
	Add(registrationform.Form) error
	SetDebugPrint(bool)
}

var DB Database

var registerPage, failurePage, succesPage []byte

var DebugPrint bool

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write(registerPage)

	case http.MethodPost:
		if DebugPrint {
			log.Println("Aplication submited.")
		}
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Form data incorrect", http.StatusBadRequest)
		}

		entry := registrationform.Form{}

		entry.Firstname = r.FormValue("vFirstName")
		entry.Lastname = r.FormValue("vLastName")

		entry.ParseEmail(r.FormValue("vEMail"))

		entry.School = r.FormValue("vSchool")
		entry.Class = r.FormValue("vClass")

		entry.ParsePhones(r.FormValue("vPhone"))
		entry.Comment = r.FormValue("vComment")

		entry.Time = time.Now()

		if ok, reason := entry.IsValid(); !ok {
			log.Printf("Failed to add entry %#v\n", entry)
			w.Header().Add("accepted", "false")
			w.Header().Add("rejection-reason", reason)
			http.Redirect(w, r, "/failure", http.StatusFound)
		} else {
			err = DB.Add(entry)
			if err != nil {
				log.Println(err)
				w.Header().Add("accepted", "false")
				w.Header().Add("rejection-reason", "internal-server-error")
				http.Redirect(w, r, "/failure", http.StatusFound)
				return
			}
			w.Header().Add("accepted", "true")
			http.Redirect(w, r, "/succes", http.StatusFound)
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func succesPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write(succesPage)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func failurePageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write(failurePage)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func initDB() {
	DB = new(db_sqlite.Database)
	DB.SetDebugPrint(DebugPrint)
}

func Start() {
	var err error
	registerPage, err = os.ReadFile("register.html")
	if err != nil {
		log.Fatal(err)
	}

	failurePage, err = os.ReadFile("failure.html")
	if err != nil {
		log.Fatal(err)
	}

	succesPage, err = os.ReadFile("succes.html")
	if err != nil {
		log.Fatal(err)
	}

	initDB()

	if err = DB.Open(); err != nil {
		log.Fatal(err)
	}
	defer DB.Close()

	if err = DB.Create(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server listening on port 8080")

	http.HandleFunc("/", registerPageHandler)
	http.HandleFunc("/succes", succesPageHandler)
	http.HandleFunc("/failure", failurePageHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
