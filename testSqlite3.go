package main

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

func mainaa() {
	db, err := sql.Open("sqlite3", "./test.db")
	//创建表
	sql_table := `
CREATE TABLE IF NOT EXISTS userinfo(
	uid INTEGER PRIMARY KEY AUTOINCREMENT,
	username VARCHAR(64) NULL,
	departname VARCHAR(64) NULL,
	created DATE NULL
);
`
	fmt.Println(db, err)
	db.Exec(sql_table)

	// insert
	stmt, err := db.Prepare("INSERT INTO userinfo(username, departname, created) values(?,?,?)")
	//checkErr(err)

	res, err := stmt.Exec("wangshubo", "国务院", "2017-04-21")
	fmt.Println(res, err)
	//checkErr(err)
	// query
	rows, err := db.Query("SELECT * FROM userinfo")
	fmt.Println(err)
	//checkErr(err)
	var uid int
	var username string
	var department string
	var created time.Time

	for rows.Next() {
		err = rows.Scan(&uid, &username, &department, &created)
		//checkErr(err)
		fmt.Println(uid, username, department, created)
	}

	rows.Close()
}
