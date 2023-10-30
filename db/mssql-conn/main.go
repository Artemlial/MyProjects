package main

import (
	// inner
	"database/sql"
	"fmt"
	"log"

	// outer
	_ "github.com/denisenkom/go-mssqldb" //driver for MS SQL Server
)

func main() {
	db, err := sql.Open("sqlserver", "sqlserver://sa:p@55w0rdMS!@localhost:1433?database=store")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	var (
		ccnum, date, cvv, exp string
		amount                float32
	)

	rows, err := db.Query("SELECT * FROM transactions;")
	if err != nil {
		log.Panicln(err)
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(&ccnum, &date, &amount, &cvv, &exp)
		if err != nil {
			log.Panicln(err)
		}
		fmt.Println(ccnum, date, amount, cvv, exp)
		if rows.Err() != nil {
			log.Panicln(err)
		}
	}
}
