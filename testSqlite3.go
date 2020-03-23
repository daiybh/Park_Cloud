package main

import (
	"database/sql"
	"fmt"
)

func maina() {
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
}
