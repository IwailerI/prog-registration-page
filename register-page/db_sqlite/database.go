package db_sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"web-server/registrationform"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	DebugPrint bool
	driver     string
	path       string
}

func (d *Database) SetDebugPrint(v bool) {
	d.DebugPrint = v
}

func (d *Database) Open() error {
	d.driver = "sqlite3"
	d.path = "./database.db"

	if d.DebugPrint {
		log.Println("Opened database:", d.path)
	}

	return nil
}

func (d *Database) Close() {
	d.driver = ""
	d.path = ""
	if d.DebugPrint {
		log.Println("Closed database")
	}
}

func (d *Database) Create() error {
	if d.DebugPrint {
		log.Println("Creating Database")
	}

	db, err := sql.Open(d.driver, d.path)
	if err != nil {
		return err
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(`
	create table if not exists RegistrationRequests (
		student_id integer not null primary key,
		firstname text,
		lastname text,
		remark text,
		class text,
		school text,
		phones text,
		email text,
		info text,
		time text default current_timestamp
	);
	`)

	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return errors.New(err.Error() + "\n" + err2.Error())
		}
		return err
	} else {
		err = tx.Commit()
		return err
	}
}

func (d *Database) Add(form registrationform.Form) error {
	if d.DebugPrint {
		log.Printf("Adding entry %#v\n", form.EscapeSQL(true))
	}

	db, err := sql.Open(d.driver, d.path)
	if err != nil {
		return err
	}
	defer db.Close()

	f := form.EscapeSQL(true)

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	_, err = tx.Exec(fmt.Sprintf(`
	insert into RegistrationRequests 
		(firstname, lastname, email, school, class, phones, info, time, remark)
		values('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '')
	`,
		f.Firstname, f.Lastname, f.Email, f.School, f.Class, f.Phones, f.Info, f.Time,
	))

	if err != nil {
		err2 := tx.Rollback()
		if err2 != nil {
			return errors.New(err.Error() + "\n" + err2.Error())
		}
		return err
	} else {
		err = tx.Commit()
		return err
	}
}
