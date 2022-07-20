package database

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"
	"web-server/registrationform"

	_ "github.com/mattn/go-sqlite3"
)

var db *sql.DB
var DebugPrint bool = false

func Open() error {
	var err error
	db, err = sql.Open("sqlite3", "./database.db")
	if DebugPrint {
		log.Println("Opened database:", err)
	}
	return err
}

func Close() {
	if db != nil {
		db.Close()
	}
	if DebugPrint {
		log.Println("Closed database")
	}
}

func Create() error {
	if DebugPrint {
		log.Println("Creating Database")
	}
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
	if DebugPrint {
		log.Printf("Added entry %#v\n", form)
	}
	_, err := db.Exec(fmt.Sprintf(`
	insert into RegistrationRequests 
		(firstname, lastname, email, school, class, phones, comment)
		values(%q, %q, %q, %q, %q, %q, %q)
	`,
		form.Firstname, form.Lastname, form.Email, form.School, form.Class, form.GetPhones(), form.Comment,
	))
	return err
}

// type Form struct {
// 	Firstname string
// 	Lastname  string
// 	Email     string
// 	School    string
// 	Class     string
// 	Phones    []Phone
// 	Comment   string
// }

func Export(filename string) error {
	fout, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer fout.Close()

	out := bufio.NewWriter(fout)
	defer out.Flush()

	rows, err := db.Query("SELECT * FROM RegistrationRequests")
	if err != nil {
		return err
	}

	var fname, lname, email, school, class, phone, comment string
	var id int64

	for rows.Next() {
		rows.Scan(&id, &fname, &lname, &email, &school, &class, &phone, &comment)
		fname = strings.ReplaceAll(fname, ",", "\\,")
		lname = strings.ReplaceAll(lname, ",", "\\,")
		email = strings.ReplaceAll(email, ",", "\\,")
		school = strings.ReplaceAll(school, ",", "\\,")
		class = strings.ReplaceAll(class, ",", "\\,")
		phone = strings.ReplaceAll(phone, ",", "\\,")
		comment = strings.ReplaceAll(comment, ",", "\\,")

		fmt.Fprintf(out, "%d,%s,%s,%s,%s,%s,%s,%s\n", id, fname, lname, email, school, class, phone, comment)
	}

	return nil
}
