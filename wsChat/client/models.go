package main


type Chat struct {
	id    int64
	name  string
	users []*User
}

type User struct {
	name   string
	tag    string
}


