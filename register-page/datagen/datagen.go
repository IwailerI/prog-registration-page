package main

import (
	"bufio"
	"database/sql"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

type EscapedForm struct {
	Firstname string
	Lastname  string
	Email     string
	School    string
	Class     string
	Phones    string
	Info      string
	Time      string
}

func main() {
	fin, err := os.Open(`D:\All\Coding\work\progmeistars-web-server\register-page\datagen\MOCK_DATA.csv`)
	if err != nil {
		panic(err)
	}
	defer fin.Close()

	in := csv.NewReader(bufio.NewReader(fin))

	db, err := sql.Open("sqlite3", "./mock_database.db")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		panic(err)
	}

	_, err = tx.Exec(`
	drop table if exists RegistrationRequests;
	create table RegistrationRequests (
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
			panic(fmt.Sprint(err, "\n", err2))
		}
		panic(err)
	}

	for {
		//id,first_name,last_name,email,school,class,phones,info,time
		record, err := in.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			panic(err)
		}

		for i := range record {
			record[i] = strings.ReplaceAll(record[i], "'", "''")
		}

		f := EscapedForm{}
		f.Firstname = record[1]
		f.Lastname = record[2]
		f.Email = record[3]
		f.School = record[4]
		f.Class = record[5]
		f.Phones = record[6]
		f.Info = record[7]
		f.Time = record[8] + ".000"

		_, err = tx.Exec(fmt.Sprintf(`
		insert into RegistrationRequests 
			(firstname, lastname, email, school, class, phones, info, time, remark)
			values('%s', '%s', '%s', '%s', '%s', '%s', '%s', '%s', '')
		`,
			f.Firstname, f.Lastname, f.Email, f.School, f.Class, f.Phones, f.Info, f.Time,
		))
		if err != nil {
			fmt.Println(f)
			err2 := tx.Rollback()
			if err2 != nil {
				panic(fmt.Sprint(err, "\n", err2))
			}
			panic(err)
		}
	}

	err = tx.Commit()
	if err != nil {
		panic(err)
	}
}
