package database

import (
	"database/sql"
	"fmt"
	"web-server/registrationform"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB

func Open() error {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	return err
}

func Close() {
	if db != nil {
		db.Close()
	}
}

func Create() error {
	_, err := db.Exec(`
	create table if not exists RegistrationRequests (
		student_id integer not null primary key,
		firstname text,
		lastname text,
		email text,
		school text,
		class text,
		phones text,
		comment text
	);
	`)
	return err
}

func Add(form registrationform.Form) error {
	_, err := db.Exec(fmt.Sprintf(`
	insert into RegistrationRequests 
		(firstname, lastname, email, school, class, phones, comment)
		values(%q, %q, %q, %q, %q, %q, %q)
	`,
		form.Firstname, form.Lastname, form.Email, form.School, form.Class, form.GetPhones(), form.Comment,
	))
	return err
}
