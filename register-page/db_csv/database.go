package db_csv

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"web-server/registrationform"
)

type Database struct {
	is_open    bool
	f          *os.File
	w          *bufio.Writer
	DebugPrint bool
}

func (d *Database) SetDebugPrint(v bool) {
	d.DebugPrint = v
}

func (d *Database) Open() error {
	fout, err := os.OpenFile("database.csv", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	if d.DebugPrint {
		log.Println("Opened database:", err)
	}
	d.f = fout
	d.w = bufio.NewWriter(fout)
	d.is_open = true
	return nil
}

func (d *Database) Close() {
	d.w.Flush()
	d.f.Close()
	d.is_open = false
	if d.DebugPrint {
		log.Println("Closed database")
	}
}

func (d *Database) Create() error {
	if d.DebugPrint {
		log.Println("Creating Database")
	}
	if d.is_open {
		return nil
	}
	if err := d.Open(); err != nil {
		return err
	}
	d.Close()
	return nil
}

func (d *Database) Add(form registrationform.Form) error {
	if d.DebugPrint {
		log.Printf("Added entry %#v\n", form)
	}
	// var fname, lname, email, school, class, phone, comment string
	// fname = strings.ReplaceAll(form.Firstname, ",", "\\,")
	// lname = strings.ReplaceAll(form.Lastname, ",", "\\,")
	// email = strings.ReplaceAll(form.Email, ",", "\\,")
	// school = strings.ReplaceAll(form.School, ",", "\\,")
	// class = strings.ReplaceAll(form.Class, ",", "\\,")
	// phone = strings.ReplaceAll(form.GetPhones(), ",", "\\,")
	// comment = strings.ReplaceAll(form.Comment, ",", "\\,")

	_, err := fmt.Fprintf(d.w, "%q,%q,%q,%q,%q,%q,%q\n",
		form.Firstname,
		form.Lastname,
		form.Email,
		form.School,
		form.Class,
		form.GetPhones(),
		form.Comment,
	)

	d.w.Flush()

	return err
}
