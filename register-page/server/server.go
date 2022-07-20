package server

import (
	"errors"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
	"web-server/db_csv"
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

var Port string
var DebugPrint bool

type DatabaseType int

const (
	SQLite DatabaseType = iota
	CSV
)

var databaseTypeString = [...]string{"sqlite", "csv"}

func (dt DatabaseType) String() string {
	if int(dt) >= len(databaseTypeString) || dt < 0 {
		return "error"
	}
	return databaseTypeString[dt]
}

func (dt *DatabaseType) Set(v string) error {
	for i, n := range databaseTypeString {
		if v == n {
			*dt = DatabaseType(i)
			return nil
		}
	}
	return errors.New("must be one of " + strings.Join(databaseTypeString[:], ", "))
}

func (dt *DatabaseType) Type() string {
	return "DatabaseType"
}

var DBType DatabaseType = SQLite

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
	switch DBType {
	case SQLite:
		DB = new(db_sqlite.Database)
	case CSV:
		DB = new(db_csv.Database)
	default:
		log.Fatal("Unsuported database type")
	}

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

	log.Println("Server listening on port " + Port)

	http.HandleFunc("/", registerPageHandler)
	http.HandleFunc("/succes", succesPageHandler)
	http.HandleFunc("/failure", failurePageHandler)
	log.Fatal(http.ListenAndServe(":"+Port, nil))
}
