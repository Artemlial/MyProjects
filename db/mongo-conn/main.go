package main

import (
	// inner
	"fmt"
	"log"

	// outer
	mgo "gopkg.in/mgo.v2"
)
type Transaction struct {
	CCnum      string `bson:"ccnum"`
	Date       string `bson:"date"`
	Amount     string `bson:"amount"`
	Cvv        string `bson:"cvv"`
	Expiration string `bson:"exp"`
}

func main(){
	session,err:=mgo.Dial("127.0.0.1")
	if err!=nil{
		log.Panicln(err)
	}
	defer session.Close()

	res:= make([]Transaction,0)

	if err:=session.DB("store").C("transactions").Find(nil).All(&res); err!=nil{
		log.Panicln(err)
	}
	for _,txn := range res{
		fmt.Println(txn.CCnum,txn.Date,txn.Amount,txn.Cvv,txn.Expiration)
	}
}
