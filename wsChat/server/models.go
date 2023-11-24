package main

import "golang.org/x/net/websocket"

type Chat struct {
	id    int64
	name  string
	users []*User
}

type User struct {
	uid    int64
	name   string
	tag    string
	active bool
	conn   *websocket.Conn
}
