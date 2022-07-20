package db_sqlite

import (
	"database/sql"
	"fmt"
	"log"
	"web-server/registrationform"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	*sql.DB
	DebugPrint bool
}

func (d *Database) SetDebugPrint(v bool) {
	d.DebugPrint = v
}

func (d *Database) Open() error {
	var err error

	d.DB, err = sql.Open("sqlite3", "./database.db")
	if d.DebugPrint {
		log.Println("Opened database:", err)
	}
	return err
}

func (d *Database) Close() {
	if d.DB != nil {
		d.DB.Close()
	}
	if d.DebugPrint {
		log.Println("Closed database")
	}
}

func (d *Database) Create() error {
	if d.DebugPrint {
		log.Println("Creating Database")
	}
	_, err := d.DB.Exec(`
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

func (d *Database) Add(form registrationform.Form) error {
	if d.DebugPrint {
		log.Printf("Added entry %#v\n", form)
	}
	_, err := d.DB.Exec(fmt.Sprintf(`
	insert into RegistrationRequests 
		(firstname, lastname, email, school, class, phones, comment)
		values(%q, %q, %q, %q, %q, %q, %q)
	`,
		form.Firstname, form.Lastname, form.Email, form.School, form.Class, form.GetPhones(), form.Comment,
	))
	return err
}
