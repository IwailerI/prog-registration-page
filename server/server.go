package server

import (
	"log"
	"net/http"
	"os"
	"web-server/database"
	"web-server/registrationform"
)

var template []byte
var Port string

func registerPageHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		w.Write(template)

	case http.MethodPut:
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
		} else {
			database.Add(entry)
			w.Header().Add("accepted", "true")
		}
		http.Redirect(w, r, "/", http.StatusFound)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func Start() {
	var err error
	template, err = os.ReadFile("template.html")
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
	log.Fatal(http.ListenAndServe(":"+Port, nil))
}
