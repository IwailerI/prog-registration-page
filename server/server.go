package server

import (
	"log"
	"net/http"
	"os"
	"web-server/database"
	"web-server/registrationform"
)

var registerPage, failurePage, succesPage []byte
var Port string

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write(registerPage)

	case http.MethodPost:
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

		if ok, reason := entry.IsValid(); !ok {
			w.Header().Add("accepted", "false")
			w.Header().Add("rejection-reason", reason)
			http.Redirect(w, r, "/failure", http.StatusFound)
		} else {
			database.Add(entry)
			log.Printf("Added entry %v\n", entry)
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

	if err = database.Open(); err != nil {
		log.Fatal(err)
	}
	defer database.Close()

	if err = database.Create(); err != nil {
		log.Fatal(err)
	}

	log.Println("Server listening on port " + Port)

	http.HandleFunc("/", registerPageHandler)
	http.HandleFunc("/succes", succesPageHandler)
	http.HandleFunc("/failure", failurePageHandler)
	log.Fatal(http.ListenAndServe(":"+Port, nil))
}
