package main

import "database/sql"

type Storage interface {
	GetUserByid(int) (*User, error)
	GetIdByCredentials(string, string) (int, error)
	CreateUser(*User) error
	UpdateUser(*User) error
}

// this one's for later

// type MongoDBStorage struct{
// 	db *sometype
// }

type SqlStorage struct {
	db *sql.DB
}

func NewSqlStorage() *SqlStorage {
	return nil
}

func (s *SqlStorage) Init() {

}

func (s *SqlStorage) GetUserByid(id int) (*User, error) {
	return nil, nil
}

func (s *SqlStorage) GetIdByCredentials(uname, pswd string) (int, error) {
	return -1, nil
}

func (s *SqlStorage) CreateUser(usr *User) error {
	return nil
}

func (s *SqlStorage) UpdateUser(usr *User) error {
	return nil
}
