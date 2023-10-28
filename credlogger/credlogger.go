package main

import (
	"os"
	// "fmt"
	"log"
	"net/http"
	"time"
)

func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		log.Printf("time::%s\n\tusername::%s\n\tpassword::%s\n\tuser-agent::%s\n\tip_address::%s\n}", time.Now().String(), r.FormValue("_user"), r.FormValue("_pass"), r.UserAgent(), r.RemoteAddr)
	}else{http.Redirect(w, r, "/", 302)}
}

func main() {
	fh, err := os.OpenFile("log.txt", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0600)
	if err != nil {
		log.Fatalln(err)
	}
	defer fh.Close()
	log.SetOutput(fh)
	log.SetPrefix("login attempt {\n\t")
	log.SetFlags(0)
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.Dir("public")))
	mux.HandleFunc("/login", login)
	log.Fatal(http.ListenAndServe(":8080", mux))
}
