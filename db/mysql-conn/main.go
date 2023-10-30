package main

import (
	// inner
	"database/sql"
	"fmt"
	"log"

	// outer
	_ "github.com/go-sql-driver/mysql" //mysql driver
)

func main() {
	db, err := sql.Open("mysql", "root:p@55w0rd!@tcp(127.0.0.1:3306)/store")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	var (
		ccnum, date, cvv, exp string
		amount                float32
	)

	rows, err := db.Query("SELECT ccnum,date,amount,cvv,exp FROM transactions")
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
	}

	if rows.Err() != nil {
		log.Panicln(err)
	}
}
