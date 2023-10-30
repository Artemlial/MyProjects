package main

import (
	// inner
	"database/sql"
	"fmt"
	"log"

	// outer
	_ "github.com/lib/pq" //postgres driver
)

func main() {
	db, err := sql.Open("postgres", "postgres://postgres:p@55w0rdPG!@localhost:3000/store?sslmode=disable")
	if err != nil {
		log.Panicln(err)
	}
	defer db.Close()

	var (
		ccnum, date, cvv, exp, amount string
	)

	rows, err := db.Query("select * from transactions")

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
